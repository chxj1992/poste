package util

import "log"

type LogLevel string

const (
	Error LogLevel = "ERROR"
	Info LogLevel = "INFO"
)

func Log(level LogLevel, format string, v ...interface{}) {
	log.Printf("[" + string(level) + "] " + format, v...)
}

func LogError(format string, v ...interface{}) {
	Log(Error, format, v...)
}

func LogInfo(format string, v ...interface{}) {
	Log(Info, format, v...)
}