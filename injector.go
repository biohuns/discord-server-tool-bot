//+build wireinject

package main

import (
	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/service/batch"
	"github.com/biohuns/discord-servertool/service/cache"
	"github.com/biohuns/discord-servertool/service/config"
	"github.com/biohuns/discord-servertool/service/discord"
	"github.com/biohuns/discord-servertool/service/gcp"
	"github.com/biohuns/discord-servertool/service/steam"
	"github.com/google/wire"
)

var superSet = wire.NewSet(
	config.ProvideService,
	cache.ProvideService,
	gcp.ProvideService,
	steam.ProvideService,
	discord.ProvideService,
	batch.ProvideService,
)

func initializeMessageService() (entity.MessageService, error) {
	wire.Build(superSet)

	return nil, nil
}

func initializeBatchService() (entity.BatchService, error) {
	wire.Build(superSet)

	return nil, nil
}
