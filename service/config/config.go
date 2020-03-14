package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"golang.org/x/xerrors"
)

const googleCredentialKey = "GOOGLE_APPLICATION_CREDENTIALS"

type (
	// Service 設定サービス
	Service struct {
		c *config
	}

	config struct {
		DiscordToken      string `json:"discord_token"`
		DiscordChannelID  string `json:"discord_channel_id"`
		DiscordBotID      string `json:"discord_bot_id"`
		GCPCredentialPath string `json:"gcp_credential_path"`
		GCPProjectID      string `json:"gcp_project_id"`
		GCPZone           string `json:"gcp_zone"`
		GCPInstanceName   string `json:"gcp_instance_name"`
		ServerAddress     string `json:"server_address"`
	}
)

var (
	serviceInstance *Service
	once            sync.Once
)

// NewService サービス返却
func ProvideService() (entity.ConfigService, error) {
	var err error

	once.Do(func() {
		filePath := flag.String("config", "config.json", "config file path")
		flag.Parse()

		var b []byte
		b, err = ioutil.ReadFile(*filePath)
		if err != nil {
			err = xerrors.Errorf("read file error: %w", err)
			return
		}

		c := new(config)
		if err = json.Unmarshal(b, c); err != nil {
			err = xerrors.Errorf("json unmarshal error: %w", err)
			return
		}

		if os.Getenv(googleCredentialKey) == "" {
			if err = os.Setenv(googleCredentialKey, c.GCPCredentialPath); err != nil {
				err = xerrors.Errorf("set env error: %w", err)
				return
			}
		}

		serviceInstance = &Service{c: c}
	})

	if serviceInstance == nil {
		err = xerrors.New("service is not provided")
	}

	if err != nil {
		return nil, xerrors.Errorf("provide service error: %w", err)
	}

	return serviceInstance, nil
}

// GetDiscordConfig Discordの設定を取得
func (s *Service) GetDiscordConfig() (token, channelID, botID string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.DiscordToken, s.c.DiscordChannelID, s.c.DiscordBotID
}

// GetGCPConfig GCPの設定を取得
func (s *Service) GetGCPConfig() (projectID, zone, instanceName string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.GCPProjectID, s.c.GCPZone, s.c.GCPInstanceName
}

// GetServerConfig サーバの設定を取得
func (s *Service) GetServerConfig() (address string) {
	if s.c == nil {
		return ""
	}
	return s.c.ServerAddress
}
