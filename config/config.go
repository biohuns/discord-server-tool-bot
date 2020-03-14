package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"golang.org/x/xerrors"
)

const googleCredentialKey = "GOOGLE_APPLICATION_CREDENTIALS"

// Config 設定
type Config struct {
	DiscordToken      string `json:"discord_token"`
	DiscordChannelID  string `json:"discord_channel_id"`
	DiscordBotID      string `json:"discord_bot_id"`
	DiscordBotName    string `json:"discord_bot_name"`
	GCPCredentialPath string `json:"gcp_credential"`
	GCPProjectID      string `json:"gcp_project_id"`
	GCPZone           string `json:"gcp_zone"`
	GCPInstanceName   string `json:"gcp_instance_name"`
}

// Open 設定ファイルを読み込む
func Open(filePath string) (*Config, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, xerrors.Errorf("read file error: %w", err)
	}

	c := new(Config)
	if err := json.Unmarshal(b, c); err != nil {
		return nil, xerrors.Errorf("json unmarshal error: %w", err)
	}

	if os.Getenv(googleCredentialKey) == "" {
		if err := os.Setenv(googleCredentialKey, c.GCPCredentialPath); err != nil {
			return nil, xerrors.Errorf("set env error: %w", err)
		}
	}

	return c, nil
}
