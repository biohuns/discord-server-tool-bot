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
	configPath := flag.String(
		"config",
		"config.json",
		"config file path",
	)
	flag.Parse()

	if err := config.Init(*configPath); err != nil {
		return xerrors.Errorf("config error: %w", err)
	}

	if err := gcp.Init(); err != nil {
		return xerrors.Errorf("GCP init error: %w", err)
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
