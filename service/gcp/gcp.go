package gcp

import (
	"context"
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
	"google.golang.org/api/compute/v1"
)

// Service GCPサービス
type Service struct {
	cache    entity.CacheService
	s        *compute.InstancesService
	project  string
	zone     string
	instance string
}

// Start インスタンス開始
func (s *Service) Start() error {
	_, err := s.s.Start(s.project, s.zone, s.instance).Do()
	if err != nil {
		return xerrors.Errorf("failed to start instance: %w", err)
	}

	return nil
}

// Stop インスタンス停止
func (s *Service) Stop() error {
	_, err := s.s.Stop(s.project, s.zone, s.instance).Do()
	if err != nil {
		return xerrors.Errorf("failed to stop instance: %w", err)
	}

	return nil
}

// GetStatus インスタンスステータスを取得する
func (s *Service) GetStatus() (*entity.InstanceStatus, error) {
	i, err := s.s.Get(s.project, s.zone, s.instance).Do()
	if err != nil {
		return nil, xerrors.Errorf("failed to get instance status: %w", err)
	}

	status := entity.InstanceStatus{Name: i.Name}

	switch i.Status {
	case "PROVISIONING", "STAGING":
		status.StatusCode = entity.InstanceStatusPending
	case "RUNNING":
		status.StatusCode = entity.InstanceStatusRunning
	case "STOPPING", "SUSPENDING":
		status.StatusCode = entity.InstanceStatusStopping
	case "STOPPED", "SUSPENDED", "TERMINATED":
		status.StatusCode = entity.InstanceStatusStopped
	default:
		status.StatusCode = entity.InstanceStatusUnknown
	}

	if v, err := s.cache.Get(entity.InstanceStatusKey); err == nil {
		if last, ok := v.(entity.InstanceStatus); ok && last.StatusCode != status.StatusCode {
			status.IsStatusChanged = true
		}
	}

	return &status, nil
}

// GetAndCacheStatus インスタンスステータスを取得しキャッシュに保存する
func (s *Service) GetAndCacheStatus() (*entity.InstanceStatus, error) {
	status, err := s.GetStatus()
	if err != nil {
		return nil, xerrors.Errorf("failed to get instance status: %w", err)
	}

	if err := s.cache.Set(entity.InstanceStatusKey, *status); err != nil {
		return nil, xerrors.Errorf("failed to set to cache: %w", err)
	}

	return status, nil
}

// GetCachedStatus キャッシュされたインスタンスステータスを取得する
func (s *Service) GetCachedStatus() (*entity.InstanceStatus, error) {
	v, err := s.cache.Get(entity.InstanceStatusKey)
	if err != nil {
		return s.GetStatus()
	}

	status, ok := v.(entity.InstanceStatus)
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
func ProvideService(conf entity.ConfigService, cache entity.CacheService) (entity.InstanceService, error) {
	var err error

	once.Do(func() {
		ctx := context.Background()
		var svc *compute.Service
		svc, err = compute.NewService(ctx)
		if err != nil {
			err = xerrors.Errorf("failed to create service: %w", err)
			return
		}

		shared = &Service{
			cache:    cache,
			s:        compute.NewInstancesService(svc),
			project:  conf.Config().GCP.Project,
			zone:     conf.Config().GCP.Zone,
			instance: conf.Config().GCP.Instance,
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
