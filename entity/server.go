package entity

import "time"

// ServerStatusKey サーバー状態のキー
const ServerStatusKey = "server_status"

// ServerStatusService サーバー状態サービス
type ServerStatusService interface {
	Status() (*ServerStatus, error)
}

// ServerStatus サーバー状態
type ServerStatus struct {
	IsOnline    bool      `json:"is_online"`
	Name        string    `json:"name"`
	PlayerCount int       `json:"player_count"`
	CheckedAt   time.Time `json:"checked_at"`

	IsStatusChanged bool          `json:"is_online_changed"`
	NobodyTime      time.Duration `json:"nobody_time"`
}
