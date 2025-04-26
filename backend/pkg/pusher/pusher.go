package pusher

import (
	"encoding/json"
	"fmt"

	"github.com/coma-toast/notifapi/backend/pkg/notification"
	"github.com/coma-toast/notifapi/internal/utils"
	pushnotifications "github.com/pusher/push-notifications-go"
)

type Pusher struct {
	InstanceID string
	SecretKey  string
	Data       *utils.DataModel
}

type MessageData struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	DeepLink string `json:"deep_link,omitempty"`
}

type AlertData struct {
	Alert    MessageData            `json:"alert,omitempty"`
	MetaData map[string]interface{} `json:"data"`
}

type NotificationData struct {
	Notification MessageData            `json:"notification,omitempty"`
	MetaData     map[string]interface{} `json:"data,omitempty"`
}

type APSData struct {
	APS AlertData
}

type Request struct {
	APNS APSData          `json:"apns,omitempty"`
	FCM  NotificationData `json:"fcm,omitempty"`
	Web  NotificationData `json:"web,omitempty"`
}

func (p Pusher) buildRequest(title, body, link string, metadata map[string]interface{}) Request {
	message := MessageData{Title: title, Body: body}
	if link != "" {
		message.DeepLink = link
	}
	request := Request{
		APNS: APSData{APS: AlertData{Alert: message, MetaData: metadata}},
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

func (p Pusher) SendMessage(payload notification.Message) (string, error) {
	return p.SendMessageFull(payload)
}

func (p Pusher) SendMessageWithLink(payload notification.Message) (string, error) {
	return p.SendMessageFull(payload)
}

func (p Pusher) SendMessageFull(payload notification.Message) (string, error) {
	beamsClient, _ := pushnotifications.New(p.InstanceID, p.SecretKey)

	request := p.buildRequest(payload.Title, payload.Body, payload.Link, payload.Metadata)
	publishRequest, err := p.convertRequest(request)
	if err != nil {
		return "", err
	}

	pubId, err := beamsClient.PublishToInterests(payload.Buckets, publishRequest)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	p.Data.AddNotification(payload)

	return pubId, nil
}

func (p Pusher) CreateInterest(name string) error {

	return nil
}
