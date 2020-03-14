package gcp

import (
	"context"

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
	// Service インスタンスサービス
	Service interface {
		Start() error
		Stop() error
		Status() (*Status, error)
	}

	// Status インスタンス情報
	Status struct {
		Name   string
		Status string
	}

	service struct {
		instance     *compute.InstancesService
		projectID    string
		zone         string
		instanceName string
	}
)

// NewService サービス生成
func NewService(projectID, zone, instanceName string) (Service, error) {
	ctx := context.Background()
	svc, err := compute.NewService(ctx)
	if err != nil {
		return nil, xerrors.Errorf("gcp error: %w", err)
	}

	return &service{
		instance:     compute.NewInstancesService(svc),
		projectID:    projectID,
		zone:         zone,
		instanceName: instanceName,
	}, nil
}

// Start インスタンス開始
func (s *service) Start() error {
	_, err := s.instance.Start(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return xerrors.Errorf("gcp error: %w", err)
	}

	return nil
}

// Stop インスタンス停止
func (s *service) Stop() error {
	_, err := s.instance.Stop(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return xerrors.Errorf("gcp error: %w", err)
	}

	return nil
}

// Status インスタンス状態確認
func (s *service) Status() (*Status, error) {
	i, err := s.instance.Get(s.projectID, s.zone, s.instanceName).Do()
	if err != nil {
		return nil, xerrors.Errorf("gcp error: %w", err)
	}

	return &Status{
		Name:   i.Name,
		Status: i.Status,
	}, nil
}
