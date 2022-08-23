package main

import (
	"fmt"
	"os"
	"time"

	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/api"
	"github.com/coma-toast/notifapi/pkg/app"
	"github.com/coma-toast/notifapi/pkg/client"
	"github.com/coma-toast/notifapi/pkg/notification"
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

	go api.RunAPI()

	//* dev code
	time.Sleep(time.Second * 5)
	client := client.Client{Target: "http://127.0.0.1:10887"}
	message := notification.Message{
		Interests: []string{"hello"},
		Title:     "Test from client",
		Body:      "This is a test from the client",
		Source:    "main.go",
	}

	response, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.Status)
	fmt.Println("message sent")
	dontExit := make(chan bool)
	// Waiting for a channel that never comes...
	<-dontExit
}
