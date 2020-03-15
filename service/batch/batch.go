package batch

import (
	"fmt"
	"sync"
	"time"

	"github.com/biohuns/discord-servertool/util"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const (
	interval        = 1 * time.Minute
	warningMinutes  = 5 * time.Minute
	shutdownMinutes = 10 * time.Minute
)

// BatchService バッチサービス
type Service struct {
	log      entity.LogService
	cache    entity.CacheService
	instance entity.InstanceService
	message  entity.MessageService
	server   entity.ServerStatusService
}

func (s Service) Start() {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			start := time.Now()
			if err := s.checkServerStatus(); err != nil {
				s.log.Error(xerrors.Errorf("failed to check server status: %w", err))
			}
			s.log.Info(
				fmt.Sprintf("finish batch (%s)", time.Now().Sub(start)),
			)
		}
	}
}

func (s Service) checkServerStatus() error {
	serverStatus, err := s.server.Status()
	if err != nil {
		return xerrors.Errorf("failed to check instance status: %w", err)
	}

	// サーバー状態変更検知

	if serverStatus.IsStatusChanged {
		_ = s.message.Send("", util.ServerStatusText(
			serverStatus.IsOnline,
			serverStatus.GameName,
			serverStatus.PlayerCount,
			serverStatus.MaxPlayerCount,
			serverStatus.Map,
		))
	}

	// 起動した状態での放置防止

	if !serverStatus.IsNobody() || serverStatus.NobodyTime < warningMinutes {
		return nil
	}

	if serverStatus.NobodyTime < shutdownMinutes {
		leftTime := shutdownMinutes - serverStatus.NobodyTime
		_ = s.message.Send("",
			fmt.Sprintf("```[%s]\n自動停止まで: %dm%ds```",
				serverStatus.GameName,
				int(leftTime.Minutes()),
				int(leftTime.Seconds())-60*int(leftTime.Minutes()),
			),
		)
		return nil
	}

	_ = s.message.Send("", fmt.Sprintf("```[%s]\n自動停止中...```", serverStatus.GameName))

	if err := s.instance.Stop(); err != nil {
		return xerrors.Errorf("failed to stop instance: %w", err)
	}

	return nil
}

var (
	shared *Service
	once   sync.Once
)

// NewService サービス返却
func ProvideService(
	log entity.LogService,
	cache entity.CacheService,
	instance entity.InstanceService,
	message entity.MessageService,
	server entity.ServerStatusService,
) (entity.BatchService, error) {
	once.Do(func() {
		shared = &Service{
			log:      log,
			cache:    cache,
			instance: instance,
			message:  message,
			server:   server,
		}
	})

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
