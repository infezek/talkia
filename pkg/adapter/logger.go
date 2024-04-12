package adapter

import (
	"encoding/json"
	"fmt"

	"github.com/infezek/app-chat/pkg/domain/adapter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/sirupsen/logrus"
)

type Logger struct {
	Params map[string]string
	Remote *newrelic.Application
	Local  *logrus.Entry
}

func NewLogger(local, remote interface{}) (*Logger, error) {
	log := Logger{
		Params: map[string]string{},
	}
	if local != nil {
		if p, ok := local.(*logrus.Entry); ok {
			log.Local = p
		} else {
			return nil, fmt.Errorf("log local is invalid")
		}
	}
	if remote != nil {
		if app, ok := remote.(*newrelic.Application); ok {
			log.Remote = app
		} else {
			return nil, fmt.Errorf("log remote is invalid")
		}
	}
	return &log, nil
}

func (l *Logger) AddParam(key, value string) {
	if key == "" || value == "" {
		return
	}
	l.Params[key] = value
	if l.Local != nil {
		l.Local = l.Local.WithFields(logrus.Fields{
			key: value,
		})
	}
	l.Local = l.Local.WithFields(logrus.Fields{key: value})
}

func (l *Logger) AddParams(params adapter.Params) {
	for key, p := range params {
		l.Params[key] = p
	}
}

func (l *Logger) Info(message string) {
	if l.Remote != nil {
		params := l.Params
		params["MESSAGE"] = message
		value, _ := json.Marshal(params)
		l.Remote.RecordLog(newrelic.LogData{
			Message:  string(value),
			Severity: "INFO",
		})
	}
	if l.Local != nil {
		l.Local.Info(message)
	}
}

func (l *Logger) Error(message string) {
	if l.Remote != nil {
		params := l.Params
		params["MESSAGE"] = message
		value, _ := json.Marshal(params)
		l.Remote.RecordLog(newrelic.LogData{
			Message:  string(value),
			Severity: "ERROR",
		})
	}
	if l.Local != nil {
		l.Local.Error(message)
	}
}
