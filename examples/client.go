package notifications

import (
	"fmt"

	"github.com/coma-toast/notifapi/pkg/client"
	"github.com/coma-toast/notifapi/pkg/notification"
)

type Notification struct {
	Target string
}

func (n *Notification) SendMessage(title, body string) error {
	client := client.Client{Target: n.Target} // the NotifAPI server (i.e. "http://alerts.mysite.com:1234")
	message := notification.Message{
		Interests: []string{"hello"},
		Title:     title,
		Body:      body,
		Source:    "Source Application Name", // what application are you sending this from
	}

	response, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(response.Status)

	return nil
}
