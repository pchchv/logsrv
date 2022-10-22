package logging

import (
	"encoding/json"
	"fmt"
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

// Logs an access entry with call duration and status code
func Access(r *http.Request, start time.Time, statusCode int) {
	e := access(r, start, statusCode, nil)
	var msg string
	if len(r.URL.RawQuery) == 0 {
		msg = fmt.Sprintf("%v ->%v %v", statusCode, r.Method, r.URL.Path)
	} else {
		msg = fmt.Sprintf("%v ->%v %v?...", statusCode, r.Method, r.URL.Path)
	}
	if statusCode >= 200 && statusCode <= 399 {
		e.Info(msg)
	} else if statusCode >= 400 && statusCode <= 499 {
		e.Warn(msg)
	} else {
		e.Error(msg)
	}
}

// Logs an error while accessing
func AccessError(r *http.Request, start time.Time, err error) {
	e := access(r, start, 0, err)
	e.Errorf("ERROR ->%v %v", r.Method, r.URL.Path)
}

func access(r *http.Request, start time.Time, statusCode int, err error) *logrus.Entry {
	url := r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}
	fields := logrus.Fields{
		"type":       "access",
		"@timestamp": start,
		"remote_ip":  getRemoteIp(r),
		"host":       r.Host,
		"url":        url,
		"method":     r.Method,
		"proto":      r.Proto,
		"duration":   time.Since(start).Nanoseconds() / 1000000,
		"User_Agent": r.Header.Get("User-Agent"),
	}
	if statusCode != 0 {
		fields["response_status"] = statusCode
	}
	if err != nil {
		fields[logrus.ErrorKey] = err.Error()
	}
	setCorrelationIds(fields, r.Header)
	cookies := map[string]string{}
	for _, c := range r.Cookies() {
		if !contains(AccessLogCookiesBlacklist, c.Name) {
			cookies[c.Name] = c.Value
		}
	}
	if len(cookies) > 0 {
		fields["cookies"] = cookies
	}

	return Logger.WithFields(fields)
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

func getRemoteIp(r *http.Request) string {
	if r.Header.Get("X-Cluster-Client-Ip") != "" {
		return r.Header.Get("X-Cluster-Client-Ip")
	}
	if r.Header.Get("X-Real-Ip") != "" {
		return r.Header.Get("X-Real-Ip")
	}
	return strings.Split(r.RemoteAddr, ":")[0]
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
