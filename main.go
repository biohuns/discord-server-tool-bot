package main

import (
	"fmt"
	"os"

	"golang.org/x/xerrors"
)

var exit = make(chan int)

func main() {
	log, err := initLogService()
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	if err := listenStart(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("message listening...")

	if err := batchStart(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
	log.Info("batch started")

	code := <-exit
	os.Exit(code)
}

func listenStart() error {
	message, err := initMessageService()
	if err != nil {
		return xerrors.Errorf("failed to init message service: %w", err)
	}

	if err := message.Start(); err != nil {
		return xerrors.Errorf("failed to start message service: %w", err)
	}

	return nil
}

func batchStart() error {
	batch, err := initBatchService()
	if err != nil {
		return xerrors.Errorf("failed to init batch service: %w", err)
	}

	go batch.Start()

	return nil
}
