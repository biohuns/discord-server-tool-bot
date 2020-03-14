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
	DiscordToken      string `json:"token"`
	DiscordChannelID  string `json:"channel_id"`
	DiscordBotID      string `json:"bot_id"`
	DiscordBotName    string `json:"bot_name"`
	GCPCredentialPath string `json:"credential"`
	GCPProjectID      string `json:"project_id"`
	GCPZone           string `json:"zone"`
	GCPInstanceName   string `json:"instance_name"`
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
		if err := os.Setenv(googleCredentialKey, c.GCP.Credential); err != nil {
			return nil, xerrors.Errorf("set env error: %w", err)
		}
	}

	return c, nil
}
