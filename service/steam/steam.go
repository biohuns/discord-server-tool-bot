package steam

import (
	"sync"
	"time"

	"github.com/Ronny95/goseq"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const timeout = 1 * time.Second

var (
	serviceInstance *Service
	once            sync.Once
)

// Service Steamサービス
type Service struct {
	cache  entity.CacheService
	server goseq.Server
}

// Status サーバ状態取得
func (s Service) Status() (*entity.ServerStatus, error) {
	last := new(entity.ServerStatus)
	if err := s.cache.Get(entity.ServerStatusKey, last); err != nil {
		last = nil
	}

	current := s.getCurrentStatus()

	if last == nil {
		if err := s.cache.Set(entity.ServerStatusKey, current); err != nil {
			return nil, err
		}
		return current, nil
	}

	if last.IsOnline != current.IsOnline {
		current.IsStatusChanged = true
	}

	if last.PlayerCount == 0 && current.PlayerCount == 0 {
		current.NobodyTime = last.NobodyTime + current.CheckedAt.Sub(last.CheckedAt)
	}

	if err := s.cache.Set(entity.ServerStatusKey, current); err != nil {
		return nil, err
	}

	return current, nil
}

func (s *Service) getCurrentStatus() *entity.ServerStatus {
	status := &entity.ServerStatus{
		CheckedAt: time.Now(),
	}

	if info, err := s.server.Info(timeout); err == nil {
		status.IsOnline = true
		status.Name = info.GetGame()
		//status.PlayerCount = int(info.GetPlayers())
	}

	return status
}

// ProvideService サービス返却
func ProvideService(config entity.ConfigService, cache entity.CacheService) (entity.ServerStatusService, error) {
	var err error

	once.Do(func() {
		address := config.GetServerConfig()

		server := goseq.NewServer()
		if err := server.SetAddress(address); err != nil {
			err = xerrors.Errorf("failed to set address: %w", err)
			return
		}

		serviceInstance = &Service{
			cache:  cache,
			server: server,
		}
	})

	if serviceInstance == nil {
		err = xerrors.Errorf("service is not provided: %w", err)
	}

	if err != nil {
		return nil, xerrors.Errorf("provide service error: %w", err)
	}

	return serviceInstance, nil
}
