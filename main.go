package main

import (
	"flag"
	"os"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/discord"
	"github.com/biohuns/discord-servertool/gcp"
	"github.com/biohuns/discord-servertool/logger"
)

var (
	stop = make(chan bool)
)

func Init() {
	configPath := flag.String(
		"config",
		"config.json",
		"config file path",
	)
	flag.Parse()

	if err := logger.Init(); err != nil {
		logger.Fatalf("logger init error: ", err)
	}

	if err := config.Init(*configPath); err != nil {
		logger.Fatalf("config init error: %s", err)
	}

	if err := os.Setenv(
		"GOOGLE_APPLICATION_CREDENTIALS",
		config.Get().GCP.Credential,
	); err != nil {
		logger.Fatalf("set env error: %s", err)
	}

	if err := gcp.Init(); err != nil {
		logger.Fatalf("gcp init error: %s", err)
	}

	if err := discord.Init(); err != nil {
		logger.Fatalf("discord init error: %s", err)
	}

	logger.Info("listening...")
	<-stop
}

func main() {
	Init()
}
