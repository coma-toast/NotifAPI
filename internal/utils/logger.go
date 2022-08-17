package utils

import (
	"os"
	"path"
	"runtime"
	"strings"

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
		l.logger.Formatter = &logrus.JSONFormatter{
			DisableTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcname := s[len(s)-1]
				_, filename := path.Split(f.File)
				return funcname, filename
			},
		}
	} else {
		l.logger.Formatter = &logrus.TextFormatter{}
	}
	l.logger.Info("Initializing Logger")
}

func (l *Logger) Info(message string) {
	l.logger.Info(message)
}
