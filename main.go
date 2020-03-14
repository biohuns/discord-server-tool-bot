package main

import (
	"fmt"
	"os"

	"github.com/biohuns/discord-servertool/logger"
	"github.com/biohuns/discord-servertool/service/config"
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
	ms, err := initializeMessageService()
	if err != nil {
		return xerrors.Errorf("listen error: %w", err)
	}

	if err := ms.Start(); err != nil {
		return xerrors.Errorf("listen error: %w", err)
	}

	return nil
}

func batchStart() error {
	panic("implement here")
}
