package app

import (
	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/notification"
)

type App struct {
	Config          utils.Config
	Data            utils.DataModel
	Logger          utils.Logger
	NotifierTargets []notification.Notifier
}

func (a *App) SendMessage(interests []string, title, body, source, link string, metadata map[string]interface{}) ([]string, []error) {
	var errors []error
	var ids []string
	for _, notifier := range a.NotifierTargets {
		id, err := notifier.SendMessageFull(interests, title, body, source, link, metadata)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		ids = append(ids, id)
	}

	return ids, errors
}
