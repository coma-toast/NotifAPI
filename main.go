package main

import (
	"flag"
	"os"

	"github.com/coma-toast/notifapi/backend/pkg/api"
	"github.com/coma-toast/notifapi/backend/pkg/app"
	"github.com/coma-toast/notifapi/backend/pkg/discord" // Add this import statement
	"github.com/coma-toast/notifapi/backend/pkg/notification"
	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/ipinfo/go/v2/ipinfo"
)

func main() {
	app := app.App{}

	configPath := flag.String("conf", ".", "Path for the config file.")
	flag.Parse()

	app.Config = *utils.GetConf(*configPath)
	app.Logger.Init(false, app.Config.LogFilePath+"notifapi.log")
	app.Logger.Info("App initialized")
	app.Data.Init(&app.Config)
	hostname, err := os.Hostname()
	if err != nil {
		app.Logger.Error(err)
	}

	app.NotifierTargets = []notification.Notifier{
		// pusher.Pusher{InstanceID: app.Config.InstanceID, SecretKey: app.Config.SecretKey, Data: &app.Data},
		discord.Discord{URL: app.Config.DiscordWebhook, Data: &app.Data},
	}
	if app.Config.Name == "" {
		app.Config.Name = hostname
	}

	app.Logger.Info("App started")
	startupMessage := notification.Message{
		Buckets: []string{"server"},
		Title:   "NotifAPI is starting up",
		Body:    "NotifAPI is starting up on " + app.Config.Name,
		Server:  app.Config.Name,
		RequestData: ipinfo.Core{
			Hostname: hostname,
		},
	}

	ids, errors := app.SendMessage(startupMessage)
	app.Logger.ProcessSendMessageResults(ids, errors)

	api := api.API{App: &app}

	go api.RunAPI()

	dontExit := make(chan bool)
	// Waiting for a channel that never comes...
	<-dontExit
}
