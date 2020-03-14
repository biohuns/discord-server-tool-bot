package main

import (
	"fmt"
	"os"

	"github.com/biohuns/discord-servertool/service/config"

	"github.com/biohuns/discord-servertool/handler"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/biohuns/discord-servertool/service/discord"
	"github.com/biohuns/discord-servertool/service/gcp"
	"github.com/google/wire"
	"golang.org/x/xerrors"
)

var superSet = wire.NewSet(
	config.NewService,
	gcp.NewService,
	discord.NewService,
)

func main() {
	if err := listenStart(); err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}

	exit := make(chan int, 1)
	<-exit
}

func listenStart() error {
	cs, err := config.NewService()
	if err != nil {
		return xerrors.Errorf("config error: %w", err)
	}

	gs, err := gcp.NewService(cs)
	if err != nil {
		return xerrors.Errorf("gcp init error: %w", err)
	}

	ds, err := discord.NewService(cs)
	if err != nil {
		return xerrors.Errorf("discord init error: %w", err)
	}

	return ds.Start(handler.NewHandler(gs, ds))
}

func batchStart() error {
	return nil
}
