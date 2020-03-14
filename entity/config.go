package entity

type (
	// ConfigService 設定サービス
	ConfigService interface {
		GetDiscordConfig() (token, channelID, botID string)
		GetGCPConfig() (projectID, zone, instanceName string)
	}
)
