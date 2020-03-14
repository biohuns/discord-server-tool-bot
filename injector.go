//+build wireinject

package main

import (
	"github.com/biohuns/discord-servertool/entity"
	"github.com/google/wire"
)

func initializeMessageService() (entity.MessageService, error) {
	wire.Build(superSet)

	return nil, nil
}
