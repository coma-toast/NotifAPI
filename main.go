package main

import (
	"os"

	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/api"
	"github.com/coma-toast/notifapi/pkg/app"
	"github.com/coma-toast/notifapi/pkg/pusher"
)

func main() {
	app := app.App{}

	app.Config = *utils.GetConf()
	app.Logger.Init(false)
	app.Logger.Info("App initialized")
	app.Data.Init(app.Config.DBFilePath)
	hostname, err := os.Hostname()
	if err != nil {
		app.Logger.Error(err)
	}
	app.Notifier = pusher.Pusher{InstanceID: app.Config.InstanceID, SecretKey: app.Config.SecretKey, Data: &app.Data}
	id, err := app.Notifier.SendMessage([]string{"hello"}, "NotifAPI", "NotifAPI is starting up on "+hostname, "main.go")
	if err != nil {
		app.Logger.ErrorWithField("Error sending message", "interest", "hello")
	}
	app.Logger.Debug(id)

	api := api.API{App: &app}

	api.RunAPI()

}
