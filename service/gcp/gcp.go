package gcp

import (
	"context"
	"fmt"
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
	"google.golang.org/api/compute/v1"
)

type (
	// Service GCPサービス
	Service struct {
		instance     *compute.InstancesService
		projectID    string
		zone         string
		instanceName string
	}

	info struct {
		Name   string
		Status entity.InstanceStatus
	}
)

var (
	serviceInstance *Service
	once            sync.Once
)

// ProvideService サービス返却
func ProvideService(cs entity.ConfigService) (entity.InstanceService, error) {
	var err error

	once.Do(func() {
		ctx := context.Background()
		var svc *compute.Service
		svc, err = compute.NewService(ctx)
		if err != nil {
			err = xerrors.Errorf("create service error: %w", err)
			return
		}

		projectID, zone, instanceName := cs.GetGCPConfig()

		serviceInstance = &Service{
			instance:     compute.NewInstancesService(svc),
			projectID:    projectID,
			zone:         zone,
			instanceName: instanceName,
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

// Start インスタンス開始
func (s *Service) Start() error {
	_, err := s.instance.Start(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return xerrors.Errorf("gcp error: %w", err)
	}

	return nil
}

// Stop インスタンス停止
func (s *Service) Stop() error {
	_, err := s.instance.Stop(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return xerrors.Errorf("gcp error: %w", err)
	}

	return nil
}

// Status インスタンス状態確認
func (s *Service) Status() (entity.InstanceInfo, error) {
	i, err := s.instance.Get(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return nil, xerrors.Errorf("gcp error: %w", err)
	}

	info := &info{Name: i.Name}

	switch i.Status {
	case "PROVISIONING":
		info.Status = entity.StatusProvisioning
	case "REPAIRING":
		info.Status = entity.StatusRepairing
	case "RUNNING":
		info.Status = entity.StatusRunning
	case "STAGING":
		info.Status = entity.StatusStaging
	case "STOPPED":
		info.Status = entity.StatusStopped
	case "STOPPING":
		info.Status = entity.StatusStopping
	case "SUSPENDED":
		info.Status = entity.StatusSuspended
	case "SUSPENDING":
		info.Status = entity.StatusSuspending
	case "TERMINATED":
		info.Status = entity.StatusTerminated
	default:
		info.Status = entity.StatusUnknown
	}

	return info, nil
}

// GetStatus インスタンス状態テキスト取得
func (s *info) GetStatus() entity.InstanceStatus {
	return s.Status
}

// GetStatus インスタンス状態テキスト取得
func (s *info) GetFormattedStatus() string {
	return fmt.Sprintf("%s: %s", s.Name, s.Status)
}
