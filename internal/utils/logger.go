package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func (l *Logger) Init(json bool) {
	l.logger = logrus.New()
	l.logger.SetReportCaller(true)
	l.logger.Out = os.Stdout
	if json {
		l.logger.Formatter = &logrus.JSONFormatter{}
		// l.logger.Formatter = &logrus.JSONFormatter{
		// 	DisableTimestamp: true,
		// 	CallerPrettyfier: func(f *runtime.Frame) (string, string) {
		// 		s := strings.Split(f.Function, ".")
		// 		funcname := s[len(s)-1]
		// 		_, filename := path.Split(f.File)
		// 		return funcname, filename
		// 	},
		// }
	} else {
		l.logger.Formatter = &logrus.TextFormatter{}
	}
	l.logger.Info("Initializing Logger")
}

func (l *Logger) Info(message string) {
	l.logger.Info(message)
}

func (l *Logger) Debug(message string) {
	l.logger.Debug(message)
}

func (l *Logger) Error(err error) {
	l.logger.Error(err)
}

func (l *Logger) ErrorWithField(message, field, value string) {
	l.logger.WithField(field, value).Error(message)
}
