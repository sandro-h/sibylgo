package main

import (
	"github.com/sirupsen/logrus"
)

// SimpleFormatter logs the message only
type SimpleFormatter struct{}

// Format renders a single log entry
func (f *SimpleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}
