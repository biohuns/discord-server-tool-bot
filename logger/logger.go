package logger

import (
	"fmt"
	"os"
	"time"
)

type LogLevel string

const (
	LogLevelFatal = "FATAL"
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
	LogLevelTrace = "TRACE"
)

var hostname string

func Init() error {
	var err error
	if hostname, err = os.Hostname(); err != nil {
		return err
	}

	return nil
}

func Fatal(message string) {
	flush(LogLevelFatal, message)
	os.Exit(1)
}
func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a...))
	os.Exit(1)
}
func Error(message string) {
	flush(LogLevelError, message)
}
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}
func Warn(message string) {
	flush(LogLevelWarn, message)
}
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}
func Info(message string) {
	flush(LogLevelInfo, message)
}
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}
func Debug(message string) {
	flush(LogLevelDebug, message)
}
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}
func Trace(message string) {
	flush(LogLevelTrace, message)
}
func Tracef(format string, a ...interface{}) {
	Trace(fmt.Sprintf(format, a...))
}
func flush(level, message string) {
	fmt.Printf(
		"%-5s %s [%s] %s\n",
		level,
		time.Now().Format("2006-01-02 15:04:05.000"),
		hostname,
		message,
	)
}
