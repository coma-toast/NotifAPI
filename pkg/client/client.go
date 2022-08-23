package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/coma-toast/notifapi/pkg/notification"
)

type Client struct {
	Target string
}

func (c *Client) SendMessage(message notification.Message) (*http.Response, error) {
	json_data, err := json.Marshal(message)
	if err != nil {
		return &http.Response{}, err
	}

	resp, err := http.Post(c.Target+"/api/notify", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
