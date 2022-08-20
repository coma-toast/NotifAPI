package notification

type Message struct {
	Interests []string               `json:"interests"`
	Title     string                 `json:"title"`
	Body      string                 `json:"body"`
	Link      string                 `json:"link"`
	Source    string                 `json:"source"`
	Metadata  map[string]interface{} `json:"metadata"`
}

type Notifier interface {
	SendMessage(interests []string, title, body, source string) (string, error)
	SendMessageWithLink(interests []string, title, body, link, source string) (string, error)
	SendMessageFull(interests []string, title, body, link, source string, metadata map[string]interface{}) (string, error)
}
