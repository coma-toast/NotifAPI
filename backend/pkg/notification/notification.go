package notification

import "github.com/ipinfo/go/v2/ipinfo"

type Message struct {
	Buckets     []string       `json:"buckets"`
	Title       string         `json:"title"`
	Body        string         `json:"body"`
	Link        string         `json:"link,omitempty"`
	Server      string         `json:"source"`
	RequestData ipinfo.Core    `json:"request_data,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type Notifier interface {
	SendMessage(message Message) (string, error)
}
