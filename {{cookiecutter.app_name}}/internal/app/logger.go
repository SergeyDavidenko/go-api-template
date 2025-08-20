package app

import (
	"time"

	"github.com/sirupsen/logrus"
)

// initLogger initializes structured logging
func initLogger() *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Add service name to all logs
	log.AddHook(&ServiceNameHook{ServiceName: "{{cookiecutter.app_name}}"})

	return log
}

// ServiceNameHook adds service name to all log entries
type ServiceNameHook struct {
	ServiceName string
}

func (h *ServiceNameHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *ServiceNameHook) Fire(entry *logrus.Entry) error {
	entry.Data["service"] = h.ServiceName
	return nil
}
