package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coma-toast/notifapi/pkg/app"
	"github.com/coma-toast/notifapi/pkg/notification"
	"github.com/gorilla/mux"
)

type API struct {
	App *app.App
}

type JSONResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
	Data  string `json:"data,omitempty"`
}

// * https://github.com/gorilla/mux#serving-single-page-applications

func (api *API) RunAPI() {
	r := mux.NewRouter()
	//  r.HandleFunc("/", api.IndexHandler)
	r.HandleFunc("/api/ping", api.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/notify", api.NotifyHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/history/{date}", api.HistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/recent/{limit}", api.RecentHandler).Methods(http.MethodGet)

	spa := spaHandler{staticPath: "./notifapi-react/build", indexPath: "./notifapi-react/public/index.html"}
	r.PathPrefix("/").Handler(spa)

	// Serve static files
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		api.App.Logger.ErrorWithField("Error getting root directory. Static files may not be served.", "error", err.Error())
	}
	cwd = cwd + "/notifapi-react/public/static"
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(cwd))))
	api.App.Logger.Error(http.ListenAndServe(fmt.Sprintf(":%s", api.App.Config.Port), r))
}

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}

func (api *API) respondWithError(w http.ResponseWriter, code int, message string) {
	api.App.Logger.ErrorWithField(message, "error from", "api")
	server, _ := os.Hostname()
	if api.App.Config.DevMode {
		id, err := api.App.Notifier.SendMessage([]string{"debug"}, "API Error Encountered", message, server)
		if err != nil {
			api.App.Logger.ErrorWithField(err.Error(), "publishID", id)
		}
	}
	api.respondWithJSON(w, code, JSONResponse{Error: message, OK: false})
}

func (api *API) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (api *API) PingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	api.App.Logger.Debug("Request sent to /api/ping")
	results, err := api.App.Notifier.SendMessage([]string{"hello"}, "NotifAPI accessed", "The endpoint /api/ping was accessed.", "PingHandler")
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusInternalServerError, "error sending notification "+err.Error())
		return
	}
	api.App.Logger.Debug(results)
	api.respondWithJSON(w, 200, "Pong\n")
}

func (api *API) HistoryHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	startDate, err := time.Parse(time.RFC3339, vars["date"])
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadRequest, "Unable to parse time string")
		return
	}

	history, err := api.App.Data.GetHistory(startDate)
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadRequest, "Unable to retrieve history")
		return
	}

	api.respondWithJSON(w, http.StatusOK, history)
}

func (api *API) RecentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	limit, err := strconv.Atoi(vars["limit"])
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadRequest, "bad limit value: "+vars["limit"])
		return
	}

	history, err := api.App.Data.GetRecentNotifications(limit)
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadRequest, "Unable to retrieve history")
		return
	}

	api.respondWithJSON(w, http.StatusOK, history)
}

func (api *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		api.App.Logger.ErrorWithField("Error getting root directory. Static files may not be served.", "error", err.Error())
	}
	cwd = cwd + "/notifapi-react/public/static"
	http.FileServer(http.Dir(cwd))

	// templateList := []string{
	// 	"./templates/head.html",
	// 	"./templates/foot.html",
	// 	"./templates/index.html",
	// }
	// templates, err := template.New("index.html").ParseFiles(templateList...)

	// templateMain := template.Must(templates, err)

	// if err := templateMain.Execute(w, nil); err != nil {
	// 	api.App.Logger.ErrorWithField(err.Error(), "API handler", "IndexHandler")
	// }
}

func (api *API) NotifyHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data notification.Message

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadGateway, "error decoding JSON")
		return
	}

	ip, err := getIP(r)
	if err != nil {
		api.App.Logger.Error(err)
	}

	api.App.Notifier.SendMessageFull(data.Interests, data.Title, data.Body, data.Link, ip, data.Metadata)

}
