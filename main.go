package main

import (
	"flag"
	"os"

	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/api"
	"github.com/coma-toast/notifapi/pkg/app"
	"github.com/coma-toast/notifapi/pkg/discord_notifier" // Add this import statement
)

func main() {
	app := app.App{}

	configPath := flag.String("conf", ".", "Path for the config file.")
	flag.Parse()

	app.Config = *utils.GetConf(*configPath)
	app.Logger.Init(false, app.Config.LogFilePath+"notifapi.log")
	app.Logger.Info("App initialized")
	app.Data.Init(app.Config.DBFilePath)
	hostname, err := os.Hostname()
	if err != nil {
		app.Logger.Error(err)
	}

	if app.Config.Name == "" {
		app.Config.Name = hostname
	}
	// app.Notifier = pusher.Pusher{InstanceID: app.Config.InstanceID, SecretKey: app.Config.SecretKey, Data: &app.Data}
	app.Notifier = discord_notifier.NewDiscordNotifier(app.Config.DiscordWebhookURL)
	// * re-enable when back online
	id, err := app.Notifier.SendMessage([]string{"hello"}, "NotifAPI", "NotifAPI is starting up on "+app.Config.Name, "main.go")
	if err != nil {
		app.Logger.ErrorWithField("Error sending message", "interest", "hello")
	}
	// id := hostname
	app.Logger.Debug(id)

	api := api.API{App: &app}

	go api.RunAPI()

	dontExit := make(chan bool)
	// Waiting for a channel that never comes...
	<-dontExit
}
