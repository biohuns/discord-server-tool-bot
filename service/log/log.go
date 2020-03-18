package log

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

// Service ログサービス
type Service struct{}

func (s *Service) Info(v interface{}) {
	flush(entity.LogLevelInfo, v)
}

func (s *Service) Warn(v interface{}) {
	flush(entity.LogLevelWarn, v)
}

func (s *Service) Error(v interface{}) {
	flush(entity.LogLevelError, v)
}

func flush(level entity.LogLevel, v interface{}) {
	var message string

	switch v := v.(type) {
	case string:
		message = v
	default:
		message = fmt.Sprintf("%+v", v)
	}

	hostname, err := os.Hostname()
	if err != nil {
		hostname = "Unknown Host"
	}

	fmt.Printf("%-5s %s [%s] %s\n", level, time.Now().Format(time.RFC3339), hostname, message)
}

var (
	shared *Service
	once   sync.Once
)

// ProvideService サービス返却
func ProvideService() (entity.LogService, error) {
	once.Do(func() {
		shared = new(Service)
	})

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
