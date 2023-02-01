package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/app"
	"github.com/coma-toast/notifapi/pkg/notification"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/kyberbits/forge/forge"
)

type API struct {
	App *app.App
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
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
	r.HandleFunc("/api/register", api.RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/api/history/{date}", api.HistoryHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/recent/{limit}", api.RecentHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/interest/{userId}/{name}", api.InterestNameHandler).Methods(http.MethodGet, http.MethodPost)
	r.HandleFunc("/api/interest/{userId}", api.InterestUserHandler).Methods(http.MethodGet)
	r.Use()

	spa := &forge.HTTPStatic{
		FileSystem: http.FS(os.DirFS("./notifapi-react/dist")),
		NotFoundHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			file, err := os.Open("./notifapi-react/dist/index.html")
			if err != nil {
				panic(err)
			}

			io.Copy(w, file)
		}),
	}
	// spa := spaHandler{staticPath: "./notifapi-react/dist", indexPath: "./notifapi-react/dist/index.html"}
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

// func authMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		claims := &Claims{}

// 		err := api.validateUserToken(claims, w, r)
// 		if err != nil {
// 			api.env.Logger.LogError("error validating user", claims.Username, err)
// 		}
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(w, r)
// 	})
// }

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
		ids, errors := api.App.SendMessage([]string{"debug"}, "API Error Encountered", message, server, "", nil)
		api.App.Logger.ProcessSendMessageResults(ids, errors)
		if len(errors) > 0 {
			for _, e := range errors {
				api.App.Logger.ErrorWithField(e.Error(), "Successful ID's", strings.Join(ids, ","))
			}
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

func (api *API) validateUserToken(claims *Claims, w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
		}

		return fmt.Errorf("Error getting cookie")
	}

	tokenString := cookie.Value

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(api.App.Config.JWTKey), nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return err
		}
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	return nil
}

func (api *API) PingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	api.App.Logger.Debug("Request sent to /api/ping")
	ids, errors := api.App.SendMessage([]string{"hello"}, "NotifAPI accessed", "The endpoint /api/ping was accessed.", "", "PingHandler", nil)
	api.App.Logger.ProcessSendMessageResults(ids, errors)
	if len(errors) > 0 {
		var errorStrings []string
		for _, e := range errors {
			errorStrings = append(errorStrings, e.Error())
		}
		api.respondWithError(w, http.StatusInternalServerError, "error sending notification "+strings.Join(errorStrings, ","))
		return
	}
	api.respondWithJSON(w, 200, "Pong\n")
}

func (api *API) InterestNameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	if r.Method == http.MethodGet {
		interests, err := api.App.Data.GetInterestsByUserAndName(vars["userId"], vars["name"])
		if err != nil {
			api.App.Logger.Error(err)
			api.respondWithError(w, http.StatusBadRequest, "Unable to get user interests")
			return
		}

		api.respondWithJSON(w, http.StatusOK, interests)
	}

	if r.Method == http.MethodPost {
		result, err := api.App.Data.InsertInterest(vars["name"], vars["webhook"], vars["userId"])
		if err != nil {
			api.App.Logger.Error(err)
			api.respondWithError(w, http.StatusBadRequest, "Unable to add webhook")
			return
		}

		api.respondWithJSON(w, http.StatusOK, result)
	}
}
func (api *API) InterestUserHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vars := mux.Vars(r)
	interests, err := api.App.Data.GetInterestsByUser(vars["userId"])
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadRequest, "Unable to get user interests")
		return
	}

	api.respondWithJSON(w, http.StatusOK, interests)
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

	api.App.SendMessage(data.Interests, data.Title, data.Body, ip, data.Link, nil)
}

func (api *API) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var data utils.User

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadGateway, "error decoding JSON")
		return
	}

	_, err = api.App.Data.AddUser(data)
	if err != nil {
		api.App.Logger.Error(err)
		api.respondWithError(w, http.StatusBadGateway, "error adding user: "+err.Error())
		return
	}

	api.respondWithJSON(w, http.StatusOK, "User added successfully")
}
