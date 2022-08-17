package pusher

import (
	"fmt"

	pushnotifications "github.com/pusher/push-notifications-go"
)

type Pusher struct {
	InstanceID string
	SecretKey  string
}

func (p Pusher) SendMessage(category, title, message string) error {
	beamsClient, _ := pushnotifications.New(p.InstanceID, p.SecretKey)

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

	return nil
}
