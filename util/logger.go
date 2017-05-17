package util

import (
	"github.com/op/go-logging"
	"os"
)

var log = logging.MustGetLogger("logger")

var f = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func LogError(format string, v ...interface{}) {
	logging.SetFormatter(f)
	log.Errorf(format, v...)
}

func LogInfo(format string, v ...interface{}) {
	logging.SetFormatter(f)
	log.Infof(format, v...)
}

func LogDebug(format string, v ...interface{}) {
	if os.Getenv("Debug") == "true" {
		logging.SetFormatter(f)
		log.Debugf(format, v...)
	}
}