package gcp

import (
	"context"

	"github.com/biohuns/discord-servertool/config"
	"golang.org/x/xerrors"
	"google.golang.org/api/compute/v1"
)

type (
	// Service サービス
	Service struct {
		is *compute.InstancesService
	}

	// Instance インスタンス情報
	Instance struct {
		Name   string
		Status string
	}
)

// Shared 共有インスタンス
var Shared *Service

// Init サービス生成
func Init() error {
	ctx := context.Background()
	s, err := compute.NewService(ctx)
	if err != nil {
		return xerrors.Errorf("GCP error: %w", err)
	}

	Shared = &Service{is: compute.NewInstancesService(s)}

	return nil
}

// Start インスタンス開始
func (s *Service) Start() error {
	_, err := s.is.Start(config.Get().GCP.ProjectID, config.Get().GCP.Zone, config.Get().GCP.InstanceName).Do()
	if err != nil {
		return xerrors.Errorf("GCP error: %w", err)
	}

	return nil
}

// Stop インスタンス停止
func (s *Service) Stop() error {
	_, err := s.is.Stop(
		config.Get().GCP.ProjectID,
		config.Get().GCP.Zone,
		config.Get().GCP.InstanceName,
	).Do()
	if err != nil {
		return xerrors.Errorf("GCP error: %w", err)
	}

	return nil
}

// Instances インスタンス一覧
func (s *Service) Instances() ([]*Instance, error) {
	list, err := s.is.List(
		config.Get().GCP.ProjectID,
		config.Get().GCP.Zone,
	).Do()
	if err != nil {
		return nil, xerrors.Errorf("GCP error: %w", err)
	}

	var instances []*Instance
	for _, instance := range list.Items {
		instances = append(instances, &Instance{
			Name:   instance.Name,
			Status: instance.Status,
		})
	}

	return instances, nil
}
