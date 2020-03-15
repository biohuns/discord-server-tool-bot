package entity

// ConfigService 設定サービス
type ConfigService interface {
	GetDiscordConfig() (token, channelID, botID string)
	GetGCPConfig() (project, zone, instance string)
	GetServerConfig() (address string)
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
	Server struct {
		Address string `json:"address"`
	} `json:"server"`
}
