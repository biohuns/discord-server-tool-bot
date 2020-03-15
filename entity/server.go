package entity

import "time"

// ServerStatusService サーバーステータスサービス
type ServerStatusService interface {
	Status() (*ServerStatus, error)
}

// ServerStatus サーバーステータス
type ServerStatus struct {
	IsOnline        bool          `json:"is_online"`
	GameName        string        `json:"game_name"`
	PlayerCount     int           `json:"player_count"`
	MaxPlayerCount  int           `json:"max_player_count"`
	Map             string        `json:"map"`
	CheckedAt       time.Time     `json:"checked_at"`
	IsStatusChanged bool          `json:"is_online_changed"`
	NobodyTime      time.Duration `json:"nobody_time"`
}

// IsNobody サーバーにだれもいない状態
func (ss *ServerStatus) IsNobody() bool {
	return ss.IsOnline && ss.PlayerCount == 0
}
