package discord_notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DiscordWebhook represents the Discord webhook configuration.
type DiscordWebhook struct {
	URL string `json:"url"`
}

// NewDiscordNotifier creates a new instance of Discord notifier.
func NewDiscordNotifier(webhookURL string) *DiscordNotifier {
	return &DiscordNotifier{
		webhookURL: webhookURL,
	}
}

// DiscordNotifier implements the Notifier interface for sending messages to Discord.
type DiscordNotifier struct {
	webhookURL string
}

// SendMessage sends a message without a link to Discord.
func (n *DiscordNotifier) SendMessage(interests []string, title, body, source string) (string, error) {
	message, err := createMessage(interests, title, body, "", source, nil)
	if err != nil {
		return "", err
	}
	resp, err := n.send(message)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// SendMessageWithLink sends a message with a link to Discord.
func (n *DiscordNotifier) SendMessageWithLink(interests []string, title, body, link, source string) (string, error) {
	message, err := createMessage(interests, title, body, link, source, nil)
	if err != nil {
		return "", err
	}
	resp, err := n.send(message)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

// SendMessageFull sends a full message including metadata to Discord.
func (n *DiscordNotifier) SendMessageFull(interests []string, title, body, link, source string, metadata map[string]interface{}) (string, error) {
	message, err := createMessage(interests, title, body, link, source, metadata)
	if err != nil {
		return "", err
	}
	resp, err := n.send(message)
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func createMessage(interests []string, title, body, link, source string, metadata map[string]interface{}) ([]byte, error) {
	data := map[string]interface{}{
		"content": fmt.Sprintf("**%s**\n%s\n*Source:* %s\n*Interests:* %v", title, body, source, interests),
	}
	if link != "" {
		data["components"] = []map[string]interface{}{
			map[string]interface{}{
				"type":     "actions",
				"elements": []map[string]interface{}{{"type": "button", "label": "View Details", "style": "primary", "url": link}},
			},
		}
	}
	if len(metadata) > 0 {
		data["footer"] = map[string]string{"text": "Metadata"}
		for k, v := range metadata {
			data[k] = v
		}
	}
	return json.Marshal(data)
}

func (n *DiscordNotifier) send(message []byte) ([]byte, error) {
	resp, err := http.Post(n.webhookURL, "application/json", bytes.NewBuffer(message))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, nil
}
