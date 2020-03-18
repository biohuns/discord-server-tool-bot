package steam

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ronny95/goseq"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const timeout = 3 * time.Second

// Service Steamサービス
type Service struct {
	cache  entity.CacheService
	server goseq.Server
}

// Status サーバ状態取得
func (s Service) Status() (*entity.ServerStatus, error) {
	last := new(entity.ServerStatus)
	v, err := s.cache.Get(entity.ServerStatusKey)
	if err != nil {
		last = nil
	}
	last, ok := v.(*entity.ServerStatus)
	if !ok {
		last = nil
	}

	current := s.getCurrentStatus()

	if last == nil {
		if err := s.cache.Set(entity.ServerStatusKey, current); err != nil {
			return nil, xerrors.Errorf("failed to set to cache: %w", err)
		}
		return current, nil
	}

	if last.IsOnline != current.IsOnline {
		current.IsStatusChanged = true
	}

	if last.IsNobody() && current.IsNobody() {
		current.NobodyTime = last.NobodyTime + current.CheckedAt.Sub(last.CheckedAt)
	}

	if err := s.cache.Set(entity.ServerStatusKey, current); err != nil {
		return nil, xerrors.Errorf("failed to set to cache: %w", err)
	}

	return current, nil
}

func (s *Service) getCurrentStatus() *entity.ServerStatus {
	status := &entity.ServerStatus{
		CheckedAt: time.Now(),
	}

	if info, err := s.server.Info(timeout); err == nil {
		status.IsOnline = true
		status.GameName = info.GetGame()
		status.PlayerCount = int(info.GetPlayers())
		status.MaxPlayerCount = int(info.GetMaxPlayers())
		status.Map = info.GetMap()
	}

	return status
}

var (
	shared *Service
	once   sync.Once
)

// ProvideService サービス返却
func ProvideService(conf entity.ConfigService, cache entity.CacheService) (entity.ServerStatusService, error) {
	var err error

	once.Do(func() {
		server := goseq.NewServer()
		if err = server.SetAddress(fmt.Sprintf(
			"%s:%d",
			conf.Config().SteamDedicatedServer.Address,
			conf.Config().SteamDedicatedServer.Port,
		)); err != nil {
			err = xerrors.Errorf("failed to set address: %w", err)
			return
		}

		shared = &Service{
			cache:  cache,
			server: server,
		}
	})

	if err != nil {
		return nil, xerrors.Errorf("failed to provide service: %w", err)
	}

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
