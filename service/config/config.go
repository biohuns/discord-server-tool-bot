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

// Service 設定サービス
type Service struct {
	c *entity.Config
}

// GetDiscordConfig Discordの設定を取得
func (s *Service) GetDiscordConfig() (token, channelID, botID string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.Discord.Token, s.c.Discord.ChannelID, s.c.Discord.BotID
}

// GetGCPConfig GCPの設定を取得
func (s *Service) GetGCPConfig() (project, zone, instance string) {
	if s.c == nil {
		return "", "", ""
	}
	return s.c.GCP.Project, s.c.GCP.Zone, s.c.GCP.Instance
}

// GetServerConfig サーバの設定を取得
func (s *Service) GetServerConfig() (address string) {
	if s.c == nil {
		return ""
	}
	return s.c.Server.Address
}

var (
	shared *Service
	once   sync.Once
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
			err = xerrors.Errorf("failed to read config file: %w", err)
			return
		}

		c := new(entity.Config)
		if err = json.Unmarshal(b, c); err != nil {
			err = xerrors.Errorf("failed to unmarshal json: %w", err)
			return
		}

		if os.Getenv(googleCredentialKey) == "" {
			if err = os.Setenv(googleCredentialKey, c.GCP.CredentialPath); err != nil {
				err = xerrors.Errorf("failed to set google application credentials: %w", err)
				return
			}
		}

		shared = &Service{c: c}
	})

	if shared == nil {
		err = xerrors.Errorf("service is not provided: %w", err)
	}

	if err != nil {
		return nil, xerrors.Errorf("failed to provide service: %w", err)
	}

	return shared, nil
}
