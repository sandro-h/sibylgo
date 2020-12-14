package main

import (
	"github.com/sandro-h/sibylgo/util"
	log "github.com/sirupsen/logrus"
	"strings"
)

// SimpleFormatter logs the message only
type SimpleFormatter struct{}

// Format renders a single log entry
func (f *SimpleFormatter) Format(entry *log.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func getConfigLogLevel(cfg *util.Config) log.Level {
	switch strings.ToLower(cfg.GetString("log_level", "info")) {
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	case "fatal":
		return log.FatalLevel
	case "panic":
		return log.PanicLevel
	default:
		return log.InfoLevel
	}
}
