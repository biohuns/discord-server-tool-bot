package batch

import (
	"fmt"
	"sync"
	"time"

	"github.com/biohuns/discord-servertool/logger"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const (
	//interval        = 10 * time.Second
	interval = 5 * time.Second
	//warningMinutes  = 20 * time.Minute
	warningMinutes = 30 * time.Second
	//shutdownMinutes = 30 * time.Minute
	shutdownMinutes = 60 * time.Second
)

var (
	serviceInstance *Service
	once            sync.Once
)

// BatchService バッチサービス
type Service struct {
	cache    entity.CacheService
	instance entity.InstanceService
	message  entity.MessageService
	server   entity.ServerStatusService
}

func (s Service) Start() {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			if err := s.checkServerStatus(); err != nil {
				logger.Error(
					fmt.Sprintf("%+v\n", xerrors.Errorf("failed to check server status: %w", err)),
				)
			}
			logger.Info(fmt.Sprintf("finish batch: %s", time.Now().Sub(start)))
		}
	}
}

func (s Service) checkServerStatus() error {
	status, err := s.server.Status()
	if err != nil {
		return xerrors.Errorf("failed to check instance status: %w", err)
	}

	// サーバー状態変更検知
	if status.IsStatusChanged {
		if status.IsOnline {
			s.message.Send(fmt.Sprintf("```[%s]\n状態変更検知: オンライン```", status.Name))
		} else {
			s.message.Send("```状態変更検知: オフライン```")
		}
	}

	if err := s.cache.Set(entity.ServerStatusKey, status); err != nil {
		return xerrors.Errorf("failed to store cache: %w", err)
	}

	// 起動した状態での放置防止
	if warningMinutes < status.NobodyTime {
		info, err := s.instance.Status()
		if err != nil {
			return xerrors.Errorf("failed to get instance info: %w", err)
		}
		if info.GetStatus() == entity.StatusRunning {
			if shutdownMinutes < status.NobodyTime {
				if err := s.instance.Stop(); err != nil {
					return xerrors.Errorf("failed to stop instance: %w", err)
				}
				s.message.Send(fmt.Sprintf("```[%s]\n自動停止中...```", status.Name))
			} else {
				s.message.Send(fmt.Sprintf("```[%s]\n自動停止まで: %s```", status.Name, shutdownMinutes-status.NobodyTime))
			}
		}
	}

	return nil
}

// NewService サービス返却
func ProvideService(
	cache entity.CacheService,
	instance entity.InstanceService,
	message entity.MessageService,
	server entity.ServerStatusService,
) (entity.BatchService, error) {
	var err error

	once.Do(func() {
		serviceInstance = &Service{
			cache:    cache,
			instance: instance,
			message:  message,
			server:   server,
		}
	})

	if serviceInstance == nil {
		err = xerrors.New("service is not provided")
	}

	if err != nil {
		return nil, xerrors.Errorf("provide service error: %w", err)
	}

	return serviceInstance, nil
}
