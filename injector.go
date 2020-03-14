//+build wireinject

package main

import (
	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/service/config"
	"github.com/biohuns/discord-servertool/service/discord"
	"github.com/biohuns/discord-servertool/service/gcp"
	"github.com/google/wire"
)

func initializeMessageService() (entity.MessageService, error) {
	wire.Build(
		config.ProvideService,
		gcp.ProvideService,
		discord.ProvideService,
	)

	return nil, nil
}
