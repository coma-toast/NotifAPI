package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

func (api *API) RunAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/", api.IndexHandler)
	r.HandleFunc("/api/ping", api.PingHandler).Methods(http.MethodGet)
	r.HandleFunc("/api/notify", api.NotifyHandler).Methods(http.MethodPost)

	// Serve static files
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		api.App.Logger.ErrorWithField("Error getting root directory. Static files may not be served.", "error", err.Error())
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(cwd+"/static/"))))
	api.App.Logger.Error(http.ListenAndServe(fmt.Sprintf(":%s", api.App.Config.Port), r))
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
	return "", fmt.Errorf("No valid ip found")
}

func (api *API) respondWithError(w http.ResponseWriter, code int, message string) {
	api.App.Logger.ErrorWithField(message, "error from", "api")
	server, _ := os.Hostname()
	id, err := api.App.Notifier.SendMessage([]string{"debug"}, "API Error Encountered", message, server)
	if err != nil {
		api.App.Logger.ErrorWithField(err.Error(), "publishID", id)
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
	api.App.Logger.Debug("Request sent to /api/ping")
	api.App.Notifier.SendMessage([]string{"hello"}, "NotifAPI accessed", "The endpoint /api/ping was accessed.", "PingHandler")
	api.respondWithJSON(w, 200, "Pong\n")
}

func (api *API) IndexHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	templateList := []string{
		"./templates/head.html",
		"./templates/foot.html",
		"./templates/index.html",
	}
	templates, err := template.New("index.html").ParseFiles(templateList...)

	templateMain := template.Must(templates, err)

	if err := templateMain.Execute(w, nil); err != nil {
		api.App.Logger.ErrorWithField(err.Error(), "API handler", "IndexHandler")
	}
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
