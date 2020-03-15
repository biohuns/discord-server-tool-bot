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

// Status インスタンス状態確認
func (s *Service) Status() (*entity.InstanceInfo, error) {
	i, err := s.s.Get(s.project, s.zone, s.instance).Do()
	if err != nil {
		return nil, xerrors.Errorf("failed to get instance: %w", err)
	}

	info := &entity.InstanceInfo{Name: i.Name}

	switch i.Status {
	case "PROVISIONING", "STAGING":
		info.Status = entity.InstanceStatusPending
	case "RUNNING":
		info.Status = entity.InstanceStatusRunning
	case "STOPPING", "SUSPENDING":
		info.Status = entity.InstanceStatusStopping
	case "STOPPED", "SUSPENDED", "TERMINATED":
		info.Status = entity.InstanceStatusStopped
	default:
		info.Status = entity.InstanceStatusUnknown
	}

	return info, nil
}

var (
	shared *Service
	once   sync.Once
)

// ProvideService サービス返却
func ProvideService(conf entity.ConfigService) (entity.InstanceService, error) {
	var err error

	once.Do(func() {
		ctx := context.Background()
		var svc *compute.Service
		svc, err = compute.NewService(ctx)
		if err != nil {
			err = xerrors.Errorf("failed to create service: %w", err)
			return
		}

		projectID, zone, instanceName := conf.GetGCPConfig()

		shared = &Service{
			s:        compute.NewInstancesService(svc),
			project:  projectID,
			zone:     zone,
			instance: instanceName,
		}
	})

	if shared == nil {
		err = xerrors.Errorf("service is not provided: %w", err)
	}

	if err != nil {
		return nil, xerrors.Errorf("failed to provide service: %w", err)
	}

	return shared, nil
}
