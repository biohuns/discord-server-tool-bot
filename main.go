package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/discord"
	"github.com/biohuns/discord-servertool/gcp"
	"github.com/biohuns/discord-servertool/logger"
	"golang.org/x/xerrors"
)

func initialize() error {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	c, err := config.Open(*configPath)
	if err != nil {
		return xerrors.Errorf("config error: %w", err)
	}

	if gcpService, err := gcp.NewService(c.GCP.ProjectID, c.GCP.Zone, c.GCP.InstanceName); err != nil {
		return xerrors.Errorf("gcp init error: %w", err)
	}

	if err := discord.Start(); err != nil {
		return xerrors.Errorf("Discord init error: %w", err)
	}

	return nil
}

func main() {
	if err := initialize(); err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	logger.Info("listening...")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	logger.Info("terminating...")
	os.Exit(0)
}
