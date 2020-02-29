package gcp

import (
	"context"

	"github.com/biohuns/discord-servertool/config"
	"google.golang.org/api/compute/v1"
)

type (
	InstanceList compute.InstanceList
)

var (
	service         *compute.Service
	instanceService *compute.InstancesService
)

func Init() (err error) {
	ctx := context.Background()
	service, err = compute.NewService(ctx)
	if err != nil {
		return err
	}
	instanceService = compute.NewInstancesService(service)

	return nil
}

func StartInstance() error {
	_, err := instanceService.Start(
		config.Get().GCP.ProjectID,
		config.Get().GCP.Zone,
		"ark-server",
	).Do()
	if err != nil {
		return err
	}

	return nil
}

func StopInstance() error {
	_, err := instanceService.Stop(
		config.Get().GCP.ProjectID,
		config.Get().GCP.Zone,
		"ark-server",
	).Do()
	if err != nil {
		return err
	}

	return nil
}

func ListInstances() (*compute.InstanceList, error) {
	list, err := instanceService.List(
		config.Get().GCP.ProjectID,
		config.Get().GCP.Zone,
	).Do()
	if err != nil {
		return nil, err
	}

	return list, nil
}
