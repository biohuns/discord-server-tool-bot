package steam

import (
	"fmt"
	"sync"
	"time"

	"github.com/Ronny95/goseq"
	"github.com/avast/retry-go"
	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const timeout = 3 * time.Second

// Service Steamサービス
type Service struct {
	cache  entity.CacheService
	server goseq.Server
}

// GetStatus サーバステータスを取得する
func (s *Service) GetStatus() (*entity.ServerStatus, error) {
	current := new(entity.ServerStatus)

	_ = retry.Do(
		func() error {
			current.CheckedAt = time.Now()

			info, err := s.server.Info(timeout)
			if err != nil {
				return err
			}

			current.IsOnline = true
			current.GameName = info.GetName()
			current.PlayerCount = int(info.GetPlayers())
			current.MaxPlayerCount = int(info.GetMaxPlayers())
			current.Map = info.GetMap()

			return nil
		},
		retry.Attempts(3),
	)

	v, err := s.cache.Get(entity.ServerStatusKey)
	if err != nil {
		return current, nil
	}

	last, ok := v.(entity.ServerStatus)
	if !ok {
		return current, nil
	}

	if last.IsOnline != current.IsOnline {
		current.IsStatusChanged = true
	}

	if last.IsNobody() && current.IsNobody() {
		current.NobodyTime = last.NobodyTime + current.CheckedAt.Sub(last.CheckedAt)
	}

	return current, nil
}

// CacheAndCacheStatus サーバステータスを取得しキャッシュに保存する
func (s *Service) GetAndCacheStatus() (*entity.ServerStatus, error) {
	status, err := s.GetStatus()
	if err != nil {
		return nil, xerrors.Errorf("failed to get server status: %w", err)
	}

	if err := s.cache.Set(entity.ServerStatusKey, *status); err != nil {
		return nil, xerrors.Errorf("failed to set to cache: %w", err)
	}

	return status, nil
}

// GetCachedStatus キャッシュされたサーバーステータス取得
func (s *Service) GetCachedStatus() (*entity.ServerStatus, error) {
	v, err := s.cache.Get(entity.ServerStatusKey)
	if err != nil {
		return s.GetStatus()
	}

	status, ok := v.(entity.ServerStatus)
	if !ok {
		return s.GetStatus()
	}

	return &status, nil
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
