package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Discord struct {
		Token     string `json:"token"`
		ChannelID string `json:"channel_id"`
		BotID     string `json:"bot_id"`
		BotName   string `json:"bot_name"`
	}
	GCP struct {
		ProjectID    string `json:"project_id"`
		Zone         string `json:"zone"`
		InstanceName string `json:"instance_name"`
	}
}

var c = new(Config)

func Init(filePath string) error {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, c); err != nil {
		return err
	}

	return nil
}

func Get() Config {
	return *c
}
