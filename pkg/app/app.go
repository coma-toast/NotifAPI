package app

import (
	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/coma-toast/notifapi/pkg/notification"
)

type App struct {
	Config   utils.Config
	Data     utils.DataModel
	Logger   utils.Logger
	Notifier notification.Notifier
}
