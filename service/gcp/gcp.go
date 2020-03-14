package gcp

import (
	"context"
	"fmt"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
	"google.golang.org/api/compute/v1"
)

const (
	StatusProvisioning = "PROVISIONING"
	StatusRepairing    = "REPAIRING"
	StatusRunning      = "RUNNING"
	StatusStaging      = "STAGING"
	StatusStopped      = "STOPPED"
	StatusStopping     = "STOPPING"
	StatusSuspended    = "SUSPENDED"
	StatusSuspending   = "SUSPENDING"
	StatusTerminated   = "TERMINATED"
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
		Status status
	}

	status string
)

// NewService サービス生成
func NewService(cs entity.ConfigService) (entity.InstanceService, error) {
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		return nil, xerrors.Errorf("gcp error: %w", err)
	}

	projectID, zone, instanceName := cs.GetGCPConfig()

	return &Service{
		instance:     compute.NewInstancesService(svc),
		projectID:    projectID,
		zone:         zone,
		instanceName: instanceName,
	}, nil
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

	return &info{
		Name:   i.Name,
		Status: status(i.Status),
	}, nil
}

// GetStatus インスタンス状態テキスト取得
func (s *info) GetStatus() string {
	return fmt.Sprintf("%s: %s", s.Name, s.Status)
}

func (s status) String() string {
	switch s {
	case StatusProvisioning:
		return "リソース割当中"
	case StatusRepairing:
		return "修復中"
	case StatusRunning:
		return "起動"
	case StatusStaging:
		return "起動準備中"
	case StatusStopped:
		return "停止"
	case StatusStopping:
		return "停止準備中"
	case StatusSuspended:
		return "休止"
	case StatusSuspending:
		return "休止準備中"
	case StatusTerminated:
		return "終了"
	default:
		return string(s)
	}
}
