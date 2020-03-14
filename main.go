package main

import (
	"fmt"
	"os"

	"github.com/biohuns/discord-servertool/logger"
	"golang.org/x/xerrors"
)

func main() {
	if err := listenStart(); err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
	logger.Info("message listening...")

	if err := batchStart(); err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		os.Exit(1)
	}
}

func listenStart() error {
	message, err := initializeMessageService()
	if err != nil {
		return xerrors.Errorf("init service error: %w", err)
	}

	if err := message.Start(); err != nil {
		return xerrors.Errorf("message listen error: %w", err)
	}

	return nil
}

func batchStart() error {
	batch, err := initializeBatchService()
	if err != nil {
		return xerrors.Errorf("init service error: %w", err)
	}

	batch.Start()

	return nil
}
