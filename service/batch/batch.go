package batch

import (
	"fmt"
	"sync"
	"time"

	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/util"
	"golang.org/x/xerrors"
)

const (
	checkInstanceStatusInterval = 30 * time.Second
	checkServerStatusInterval   = 1 * time.Minute
	warningMinutes              = 10 * time.Minute
	shutdownMinutes             = 20 * time.Minute
)

// BatchService バッチサービス
type Service struct {
	log      entity.LogService
	instance entity.InstanceService
	server   entity.ServerStatusService
	message  entity.MessageService
}

// Start バッチ処理を開始
func (s *Service) Start() {
	isTicker := time.NewTicker(checkInstanceStatusInterval)
	defer isTicker.Stop()

	ssTicker := time.NewTicker(checkServerStatusInterval)
	defer ssTicker.Stop()

	for {
		select {
		case <-isTicker.C:
			go s.checkInstanceStatus()
		case <-ssTicker.C:
			go s.checkServerStatus()
		}
	}
}

func (s *Service) checkInstanceStatus() {
	start := time.Now()
	defer func() {
		s.log.Info(fmt.Sprintf("finish instance status batch (in %s)", time.Now().Sub(start)))
	}()

	status, err := s.instance.GetAndCacheStatus()
	if err != nil {
		s.log.Error(xerrors.Errorf("failed to check instance status: %w", err))
		return
	}

	if status.IsStatusChanged {
		_ = s.message.Send("", util.InstanceStatusText(status.Name, status.StatusCode.String()))
	}
}

func (s *Service) checkServerStatus() {
	start := time.Now()
	defer func() {
		s.log.Info(fmt.Sprintf("finish server status batch (in %s)", time.Now().Sub(start)))
	}()

	status, err := s.server.GetAndCacheStatus()
	if err != nil {
		s.log.Error(xerrors.Errorf("failed to check server status: %w", err))
		return
	}

	// サーバーステータス変更検知
	if status.IsStatusChanged {
		_ = s.message.Send("", util.ServerStatusText(
			status.IsOnline,
			status.GameName,
			status.PlayerCount,
			status.MaxPlayerCount,
			status.Map,
		))
	}

	// 起動した状態での放置防止
	if !status.IsNobody() || status.NobodyTime < warningMinutes {
		return
	}

	if status.NobodyTime < shutdownMinutes {
		leftTime := shutdownMinutes - status.NobodyTime
		_ = s.message.Send("",
			fmt.Sprintf("```[%s]\nAuto Stop Instance (In %dm%ds)```",
				status.GameName,
				int(leftTime.Minutes()),
				int(leftTime.Seconds())-60*int(leftTime.Minutes()),
			),
		)
		return
	}

	_ = s.message.Send("", fmt.Sprintf("```[%s]\nStopping Instance Automatically...```", status.GameName))

	if err := s.instance.Stop(); err != nil {
		s.log.Error(xerrors.Errorf("failed to stop instance: %w", err))
		return
	}
}

var (
	shared *Service
	once   sync.Once
)

// NewService サービス返却
func ProvideService(
	log entity.LogService,
	instance entity.InstanceService,
	server entity.ServerStatusService,
	message entity.MessageService,
) (entity.BatchService, error) {
	once.Do(func() {
		shared = &Service{
			log:      log,
			instance: instance,
			server:   server,
			message:  message,
		}
	})

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
