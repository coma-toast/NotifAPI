package app

import (
	"github.com/coma-toast/notifapi/backend/pkg/notification"
	"github.com/coma-toast/notifapi/internal/utils"
)

type App struct {
	Config          utils.Config
	Data            utils.DataModel
	Logger          utils.Logger
	NotifierTargets []notification.Notifier
}

func (a *App) SendMessage(payload notification.Message) ([]string, []error) {
	var errors []error
	var ids []string
	_, err := a.Data.AddNotification(payload)
	if err != nil {
		a.Logger.ErrorWithField("error adding notification to db", payload.Title, err.Error())
		return nil, []error{err}
	}

	for _, notifier := range a.NotifierTargets {
		a.Logger.LogMessage(payload)
		id, err := notifier.SendMessage(payload)
		if err != nil {
			a.Logger.ErrorWithField("error sending message", payload.Title, err.Error())
			errors = append(errors, err)
			continue
		}
		ids = append(ids, id)
	}

	return ids, errors
}
