package internal

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type LogFormatter struct{}

func (L *LogFormatter) Format(entry *log.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006/01/02 15:04:05")
	msg := fmt.Sprintf("%s [%s] %s\n", timestamp, strings.ToUpper(entry.Level.String()), entry.Message)
	return []byte(msg), nil
}
