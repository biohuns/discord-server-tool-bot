package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/biohuns/discord-servertool/adapter/discord"
	"github.com/biohuns/discord-servertool/adapter/gcp"
	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/handler"
	"github.com/biohuns/discord-servertool/logger"
	"golang.org/x/xerrors"
)

func start() error {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	c, err := config.Open(*configPath)
	if err != nil {
		return xerrors.Errorf("config error: %w", err)
	}

	gs, err := gcp.NewService(c.GCPProjectID, c.GCPZone, c.GCPInstanceName)
	if err != nil {
		return xerrors.Errorf("gcp init error: %w", err)
	}

	ds, err := discord.NewService(c.DiscordToken, c.DiscordChannelID, c.DiscordBotID)
	if err != nil {
		return xerrors.Errorf("Discord init error: %w", err)
	}

	handlerFunc := handler.NewHandler(gs, ds)
	return ds.Start(handlerFunc)
}

func main() {
	if err := start(); err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
	logger.Info("terminating...")
	os.Exit(0)
}
