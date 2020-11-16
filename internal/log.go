package internal

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (L *LogFormatter) Format(entry *log.Entry) ([]byte, error) {

	var fieldSlice []string
	for key, val := range entry.Data {
		fieldSlice = append(fieldSlice, fmt.Sprintf("%s=%s", key, val))
	}

	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message, strings.Join(fieldSlice, " "))
	return []byte(msg), nil
}

// log hook
type ServiceNameHook struct {
	ServiceName string
}

func (s *ServiceNameHook) Levels() []log.Level {
	return log.AllLevels
}

func (s *ServiceNameHook) Fire(entry *log.Entry) error {
	entry.Data["serviceName"] = s.ServiceName
	return nil
}
