package main

import (
	"fmt"

	"github.com/coma-toast/notifapi/utils"
	pushnotifications "github.com/pusher/push-notifications-go"
)

func main() {
	config := utils.GetConf()
	beamsClient, _ := pushnotifications.New(config.InstanceID, config.SecretKey)

	publishRequest := map[string]interface{}{
		"apns": map[string]interface{}{
			"aps": map[string]interface{}{
				"alert": map[string]interface{}{
					"title": "Hello",
					"body":  "Hello, world",
				},
			},
		},
		"fcm": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
		"web": map[string]interface{}{
			"notification": map[string]interface{}{
				"title": "Hello",
				"body":  "Hello, world",
			},
		},
	}

	pubId, err := beamsClient.PublishToInterests([]string{"hello"}, publishRequest)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Publish Id:", pubId)
	}
}
