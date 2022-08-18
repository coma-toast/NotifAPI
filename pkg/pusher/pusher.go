package pusher

import (
	"encoding/json"
	"fmt"

	pushnotifications "github.com/pusher/push-notifications-go"
)

type Pusher struct {
	InstanceID string
	SecretKey  string
}

type MessageData struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	DeepLink string `json:"deep_link,omitempty"`
}

type NotificationData struct {
	Alert        MessageData            `json:"alert,omitempty"`
	Notification MessageData            `json:"notification,omitempty"`
	MetaData     map[string]interface{} `json:"data"`
}

type APSData struct {
	APS NotificationData
}

type Request struct {
	APNS APSData          `json:"apns,omitempty"`
	FCM  NotificationData `json:"fcm,omitempty"`
	Web  NotificationData `json:"web,omitempty"`
}

func (p Pusher) buildRequest(title, body, link string, metadata map[string]interface{}) Request {
	message := MessageData{Title: title, Body: body, DeepLink: link}
	request := Request{
		APNS: APSData{APS: NotificationData{Alert: message, MetaData: metadata}},
		FCM:  NotificationData{Notification: message, MetaData: metadata},
		Web:  NotificationData{Notification: message, MetaData: metadata},
	}

	return request
}

// Convert a Request to map[string]interface{} to satisfy the beamsClient.PublishToInterests
func (p Pusher) convertRequest(request Request) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	jsonData, err := json.Marshal(request)
	if err != nil {
		return m, err
	}

	json.Unmarshal(jsonData, &m)

	return m, nil
}

func (p Pusher) SendMessage(interests []string, title, body, link string, metadata map[string]interface{}) error {
	beamsClient, _ := pushnotifications.New(p.InstanceID, p.SecretKey)

	request := p.buildRequest(title, body, link, metadata)
	publishRequest, err := p.convertRequest(request)
	if err != nil {
		return err
	}

	pubId, err := beamsClient.PublishToInterests(interests, publishRequest)
	if err != nil {
		fmt.Println(err)
		return err
	} else {
		fmt.Println("Publish Id:", pubId)
	}

	return nil
}

func (p Pusher) CreateInterest(name string) error {

	return nil
}
