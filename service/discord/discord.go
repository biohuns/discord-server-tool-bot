package discord

import (
	"fmt"
	"strings"

	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

type (
	service struct {
		session   *discordgo.Session
		channelID string
		botID     string
	}
)

// NewService サービス生成
func NewService(cs entity.ConfigService) (entity.MessageService, error) {
	session, err := discordgo.New()
	if err != nil {
		return nil, xerrors.Errorf("create session error: %w", err)
	}

	token, channelID, botID := cs.GetDiscordConfig()
	session.Token = fmt.Sprintf("Bot %s", token)

	return &service{
		session:   session,
		channelID: channelID,
		botID:     botID,
	}, nil
}

// Start ハンドラを追加して監視を開始
func (s *service) Start(handler ...interface{}) error {
	for _, h := range handler {
		handlerFunc, ok := h.(func(s *discordgo.Session, m *discordgo.MessageCreate))
		if !ok {
			return xerrors.New("handler func convert err")
		}
		s.session.AddHandler(handlerFunc)
	}

	if err := s.session.Open(); err != nil {
		return xerrors.Errorf("session open error: %w", err)
	}

	logger.Info("listening...")

	return nil
}

// Send メッセージ送信
func (s *service) Send(msg string) {
	if _, err := s.session.ChannelMessageSend(s.channelID, msg); err != nil {
		logger.Error(
			fmt.Sprintf("%+v", xerrors.Errorf("message send error: %w", err)),
		)
	}
}

// SendTo メッセージ送信（対象を取る）
func (s *service) SendTo(userID, msg string) {
	s.Send(fmt.Sprintf("<@!%s>\n%s", userID, msg))
}

// GetCommand メッセージからコマンド部分を抽出
func (s *service) GetCommand(m *discordgo.MessageCreate) string {
	cmd := strings.TrimSpace(m.Content)

	if strings.HasPrefix(cmd, fmt.Sprintf("<@%s>", s.botID)) {
		cmd = strings.Replace(cmd, fmt.Sprintf("<@%s>", s.botID), "", 1)
	} else if strings.HasPrefix(cmd, fmt.Sprintf("<@!%s>", s.botID)) {
		cmd = strings.Replace(cmd, fmt.Sprintf("<@!%s>", s.botID), "", 1)
	} else {
		return ""
	}

	return strings.TrimSpace(cmd)
}

// IsCommand 自身に向けられたコマンドであるか
func (s *service) IsCommand(m *discordgo.MessageCreate) bool {
	return s.botID != m.Author.ID &&
		m.ChannelID == s.channelID &&
		(strings.HasPrefix(m.Content, fmt.Sprintf("<@%s>", s.botID)) ||
			strings.HasPrefix(m.Content, fmt.Sprintf("<@!%s>", s.botID))) &&
		s.GetCommand(m) != ""
}
