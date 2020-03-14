package entity

import "github.com/bwmarrin/discordgo"

// MessageService メッセージサービス
type MessageService interface {
	Start(handler ...interface{}) error
	Send(msg string)
	SendTo(userID, msg string)
	GetCommand(m *discordgo.MessageCreate) string
	IsCommand(m *discordgo.MessageCreate) bool
}
