package entity

// MessageService メッセージサービス
type MessageService interface {
	Start() error
	Send(msg string)
	SendTo(userID, msg string)
}
