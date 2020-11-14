package utils

import (
	"fmt"
	"path"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func CreateRotatelogs(log_prefix, serviceName, logpath string) *rotatelogs.RotateLogs {
	log_f := fmt.Sprintf("%s_%s", log_prefix, serviceName)
	log_f = path.Join(logpath, log_f)
	rl, err := rotatelogs.New(
		log_f+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(log_f),
		rotatelogs.WithMaxAge(24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Fatalf("failed to create rotatelogs: %s", err)
	}
	return rl
}
