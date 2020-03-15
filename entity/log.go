package entity

// LogService ログサービス
type LogService interface {
	Info(v interface{})
	Warn(v interface{})
	Error(v interface{})
}

type LogLevel int

const (
	LogLevelUnknown LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

func (l LogLevel) String() string {
	switch l {
	case LogLevelUnknown:
		return "UNKNOWN"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}
