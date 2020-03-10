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
	Discord struct {
		Token     string `json:"token"`
		ChannelID string `json:"channel_id"`
		BotID     string `json:"bot_id"`
		BotName   string `json:"bot_name"`
	}
	GCP struct {
		Credential   string `json:"credential"`
		ProjectID    string `json:"project_id"`
		Zone         string `json:"zone"`
		InstanceName string `json:"instance_name"`
	}
}

var c = new(Config)

// Init 設定を初期化する
func Init(filePath string) error {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return xerrors.Errorf("read file error: %w", err)
	}

	if err := json.Unmarshal(b, c); err != nil {
		return xerrors.Errorf("json unmarshal error: %w", err)
	}

	if os.Getenv(googleCredentialKey) == "" {
		if err := os.Setenv(googleCredentialKey, c.GCP.Credential); err != nil {
			return xerrors.Errorf("set env error: %w", err)
		}
	}

	return nil
}

// Get 設定を取得する
func Get() Config {
	if c == nil {
		return Config{}
	}
	return *c
}
