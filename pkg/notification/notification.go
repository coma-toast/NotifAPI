package notification

type Notifier interface {
	SendMessage(category, title, message string) error
}
