package util

import "log"

var Debugging = false

type LogLevel string

const (
	Error LogLevel = "ERROR"
	Info LogLevel = "INFO"
	Debug LogLevel = "DEBUG"
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

func LogDebug(format string, v ...interface{}) {
	if Debugging {
		Log(Debug, format, v...)
	}
}