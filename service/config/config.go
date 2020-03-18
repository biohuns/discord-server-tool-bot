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

// Config 設定を返却する
func (s *Service) Config() entity.Config {
	if s.c == nil {
		return entity.Config{}
	}
	return *s.c
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

		if err = c.Validate(); err != nil {
			err = xerrors.Errorf("invalid config: %w", err)
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

	if err != nil {
		return nil, xerrors.Errorf("failed to provide service: %w", err)
	}

	if shared == nil {
		return nil, xerrors.New("service is not provided")
	}

	return shared, nil
}
