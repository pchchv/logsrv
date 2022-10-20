package logging

import (
	"time"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Set("info", false)
}

// Set creates a new Logger with the matching specification
func Set(level string, textLogging bool) error {
	l, err := logrus.ParseLevel(level)
	if err != nil {
		return err
	}

	logger := logrus.New()
	if textLogging {
		logger.Formatter = &logrus.TextFormatter{}
	} else {
		logger.Formatter = &LogstashFormatter{TimestampFormat: time.RFC3339Nano}
	}
	logger.Level = l
	Logger = logger
	return nil
}
