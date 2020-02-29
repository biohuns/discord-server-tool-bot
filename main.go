package main

import (
	"flag"
	"log"

	"github.com/biohuns/discord-server-tool-bot/config"
	"github.com/biohuns/discord-server-tool-bot/discord"
	"github.com/biohuns/discord-server-tool-bot/gcp"
)

var (
	stop = make(chan bool)
)

func main() {
	configPath := ""
	flag.StringVar(&configPath, "config", "config.json", "config file path")
	flag.Parse()
	if err := config.Init("config.json"); err != nil {
		log.Fatalln(err)
	}

	if err := gcp.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := discord.Init(); err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening...")
	<-stop //プログラムが終了しないようロック
	return
}
