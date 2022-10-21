package logging

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
	// The of cookies which should not be logged
	AccessLogCookiesBlacklist = []string{}
	LifecycleEnvVars          = []string{"BUILD_NUMBER", "BUILD_HASH", "BUILD_DATE"}
)

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

// Return a log entry for application logs, prefilled with the correlation ids out of the supplied request
func Application(h http.Header) *logrus.Entry {
	fields := logrus.Fields{
		"type": "application",
	}
	setCorrelationIds(fields, h)
	return Logger.WithFields(fields)
}

// Logs the start of an application with the configuration struct or map as paramter
func LifecycleStart(appName string, args interface{}) {
	fields := logrus.Fields{}
	jsonString, err := json.Marshal(args)
	if err == nil {
		err := json.Unmarshal(jsonString, &fields)
		if err != nil {
			fields["parse_error"] = err.Error()
		}
	}
	fields["type"] = "lifecycle"
	fields["event"] = "start"
	for _, env := range LifecycleEnvVars {
		if os.Getenv(env) != "" {
			fields[strings.ToLower(env)] = os.Getenv(env)
		}
	}
	Logger.WithFields(fields).Infof("starting application: %v", appName)
}

// Logs the stop of an application
func LifecycleStop(appName string, signal os.Signal, err error) {
	fields := logrus.Fields{
		"type":  "lifecycle",
		"event": "stop",
	}
	if signal != nil {
		fields["signal"] = signal.String()
	}
	if os.Getenv("BUILD_NUMBER") != "" {
		fields["build_number"] = os.Getenv("BUILD_NUMBER")
	}
	if err != nil {
		Logger.WithFields(fields).
			WithError(err).
			Errorf("stopping application: %v (%v)", appName, err)
	} else {
		Logger.WithFields(fields).Infof("stopping application: %v (%v)", appName, signal)
	}
}

// Logs the stop of an application
func ServerClosed(appName string) {
	fields := logrus.Fields{
		"type":  "application",
		"event": "stop",
	}
	if os.Getenv("BUILD_NUMBER") != "" {
		fields["build_number"] = os.Getenv("BUILD_NUMBER")
	}
	Logger.WithFields(fields).Infof("http server was closed: %v", appName)
}

func setCorrelationIds(fields logrus.Fields, h http.Header) {
	correlationId := GetCorrelationId(h)
	if correlationId != "" {
		fields["correlation_id"] = correlationId
	}

	userCorrelationId := GetUserCorrelationId(h)
	if userCorrelationId != "" {
		fields["user_correlation_id"] = userCorrelationId
	}
}
