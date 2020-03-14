package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const googleCredentialKey = "GOOGLE_APPLICATION_CREDENTIALS"

type (
	service struct {
		c *config
	}

	config struct {
		DiscordToken      string `json:"discord_token"`
		DiscordChannelID  string `json:"discord_channel_id"`
		DiscordBotID      string `json:"discord_bot_id"`
		GCPCredentialPath string `json:"gcp_credential"`
		GCPProjectID      string `json:"gcp_project_id"`
		GCPZone           string `json:"gcp_zone"`
		GCPInstanceName   string `json:"gcp_instance_name"`
	}
)

// NewService サービス生成
func NewService() (entity.ConfigService, error) {
	filePath := flag.String("config", "config.json", "config file path")
	flag.Parse()

	b, err := ioutil.ReadFile(*filePath)
	if err != nil {
		return nil, xerrors.Errorf("read file error: %w", err)
	}

	c := new(config)
	if err := json.Unmarshal(b, c); err != nil {
		return nil, xerrors.Errorf("json unmarshal error: %w", err)
	}

	if os.Getenv(googleCredentialKey) == "" {
		if err := os.Setenv(googleCredentialKey, c.GCPCredentialPath); err != nil {
			return nil, xerrors.Errorf("set env error: %w", err)
		}
	}

	return &service{c: c}, nil
}

// GetDiscordConfig Discordの設定を取得
func (s service) GetDiscordConfig() (token, channelID, botID string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.DiscordToken, s.c.DiscordChannelID, s.c.DiscordBotID
}

// GetGCPConfig GCPの設定を取得
func (s service) GetGCPConfig() (projectID, zone, instanceName string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.GCPProjectID, s.c.GCPZone, s.c.GCPInstanceName
}
