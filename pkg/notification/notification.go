package notification

type Notifier interface {
	SendMessage(interests []string, title, message, link string, metadata map[string]interface{}) error
}
