package main

import (
	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/notification"
	"github.com/coma-toast/notifapi/pkg/pusher"
)

type App struct {
	Config   utils.Config
	Data     utils.DataModel
	Logger   utils.Logger
	Notifier notification.Notifier
}

func main() {
	app := App{}

	app.Config = *utils.GetConf()
	app.Logger.Init(false)
	app.Notifier = pusher.Pusher{InstanceID: app.Config.InstanceID, SecretKey: app.Config.SecretKey}
	app.Logger.Info("App initialized")
	app.Data.Init(app.Config.DBFilePath)
	data := map[string]interface{}{
		"test data": "asdf",
	}
	app.Notifier.SendMessage([]string{"hello"}, "Test notification", "This is a test", "", data)
	// go run API
	// go run HTTP

}
