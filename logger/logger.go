package logger

import (
	"fmt"
	"os"
	"time"
)

type LogLevel string

const (
	LogLevelError = "ERROR"
	LogLevelWarn  = "WARN"
	LogLevelInfo  = "INFO"
	LogLevelDebug = "DEBUG"
)

// Error エラー出力
func Error(message string) {
	flush(LogLevelError, message)
}

// Warning 警告出力
func Warning(message string) {
	flush(LogLevelWarn, message)
}

// Info 情報出力
func Info(message string) {
	flush(LogLevelInfo, message)
}

// Debug デバッグ出力
func Debug(message string) {
	flush(LogLevelDebug, message)
}

func flush(level, message string) {
	var (
		hostname string
		err      error
	)

	if hostname, err = os.Hostname(); err != nil {
		fmt.Printf("%-5s %s %s\n", level, time.Now().Format(time.RFC3339), message)
	}

	fmt.Printf("%-5s %s [%s] %s\n", level, time.Now().Format(time.RFC3339), hostname, message)
}
