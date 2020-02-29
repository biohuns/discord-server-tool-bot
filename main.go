package main

import (
	"flag"
	"log"
	"os"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/discord"
	"github.com/biohuns/discord-servertool/gcp"
)

var (
	stop = make(chan bool)
)

func main() {
	configPath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	if err := config.Init(*configPath); err != nil {
		log.Fatalln(err)
	}

	if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.Get().GCP.Credential); err != nil {
		log.Fatalln(err)
	}

	if err := gcp.Init(); err != nil {
		log.Fatalln(err)
	}

	if err := discord.Init(); err != nil {
		log.Fatalln(err)
	}

	log.Println("Listening...")
	<-stop
	return
}
