package logging

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Generates json in logstash format
type LogstashFormatter struct {
	// If not empty use for logstash type field
	Type string
	// Sets the format used for timestamps
	TimestampFormat string
}

func (f *LogstashFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return f.FormatWithPrefix(entry, "")
}

func (f *LogstashFormatter) FormatWithPrefix(entry *logrus.Entry, prefix string) ([]byte, error) {
	fields := make(logrus.Fields)
	for k, v := range entry.Data {
		// Remove the prefix when sending the fields to logstash
		if prefix != "" && strings.HasPrefix(k, prefix) {
			k = strings.TrimPrefix(k, prefix)
		}
		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			fields[k] = v.Error()
		default:
			fields[k] = v
		}
	}
	fields["@version"] = "1"
	timeStampFormat := f.TimestampFormat
	if timeStampFormat == "" {
		timeStampFormat = time.RFC3339
	}
	fields["@timestamp"] = entry.Time.Format(timeStampFormat)
	// Set message field
	v, ok := entry.Data["message"]
	if ok {
		fields["fields.message"] = v
	}
	fields["message"] = entry.Message
	// Set level field
	v, ok = entry.Data["level"]
	if ok {
		fields["fields.level"] = v
	}
	fields["level"] = entry.Level.String()
	// Set type field
	if f.Type != "" {
		v, ok = entry.Data["type"]
		if ok {
			fields["fields.type"] = v
		}
		fields["type"] = f.Type
	}
	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
