package entity

import "golang.org/x/xerrors"

// ConfigService 設定サービス
type ConfigService interface {
	Config() Config
}

// Config 設定
type Config struct {
	Discord struct {
		Token     string `json:"token"`
		ChannelID string `json:"channel_id"`
		BotID     string `json:"bot_id"`
	} `json:"discord"`
	GCP struct {
		CredentialPath string `json:"credential_path"`
		Project        string `json:"project"`
		Zone           string `json:"zone"`
		Instance       string `json:"instance"`
	} `json:"gcp"`
	SteamDedicatedServer struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	} `json:"steam_dedicated_server"`
}

// Validate バリデーションを行う
func (c *Config) Validate() error {
	if c.Discord.Token == "" ||
		c.Discord.ChannelID == "" ||
		c.Discord.BotID == "" {
		return xerrors.Errorf("not enough discord config: %+v", c.Discord)
	}

	if c.GCP.CredentialPath == "" ||
		c.GCP.Project == "" ||
		c.GCP.Zone == "" ||
		c.GCP.Instance == "" {
		return xerrors.Errorf("not enough gcp config: %+v", c.GCP)
	}

	if c.SteamDedicatedServer.Address == "" ||
		c.SteamDedicatedServer.Port <= 0 {
		return xerrors.Errorf("not enough steam dedicated server config: %+v", c.GCP)
	}

	return nil
}
