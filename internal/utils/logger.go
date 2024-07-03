package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
	file   *os.File
}

// Write logs to multiple writers
type multiLogWriter struct {
	writers []io.Writer
}

func (m multiLogWriter) Write(p []byte) (int, error) {
	n := len(p)
	for _, w := range m.writers {
		nw, err := w.Write(p[:n])
		if err != nil {
			return nw, err
		}
		n -= nw
	}
	return n, nil
}

// Add a method to set the writers
func (m *multiLogWriter) SetWriters(writers []io.Writer) {
	m.writers = writers
}

func (l *Logger) Init(json bool, logFile string) {
	var err error
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
	// Create a custom formatter that writes to both stdout and a file
	l.file, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	multiLogWriter := multiLogWriter{}
	multiLogWriter.SetWriters([]io.Writer{os.Stdout, l.file})
	l.logger.Out = multiLogWriter
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
