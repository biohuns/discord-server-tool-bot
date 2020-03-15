package entity

// MessageService メッセージサービス
type MessageService interface {
	Start() error
	Send(userID, message string) error
}

// StatusMessage ステータスメッセージ
type StatusMessage struct {
	InstanceName   string
	InstanceStatus string
	IsOnline       bool
	GameName       string
	PlayerCount    int
}
