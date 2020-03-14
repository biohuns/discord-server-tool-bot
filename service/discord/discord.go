package discord

import (
	"fmt"
	"strings"
	"sync"

	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

// Service Discordサービス
type Service struct {
	is        entity.InstanceService
	session   *discordgo.Session
	channelID string
	botID     string
}

var (
	serviceInstance *Service
	once            sync.Once
)

// ProvideService サービス返却
func ProvideService(cs entity.ConfigService, is entity.InstanceService) (entity.MessageService, error) {
	var err error

	once.Do(func() {
		var session *discordgo.Session
		session, err = discordgo.New()
		if err != nil {
			err = xerrors.Errorf("create session error: %w", err)
			return
		}

		token, channelID, botID := cs.GetDiscordConfig()
		session.Token = fmt.Sprintf("Bot %s", token)

		serviceInstance = &Service{
			is:        is,
			session:   session,
			channelID: channelID,
			botID:     botID,
		}
	})

	if serviceInstance == nil {
		err = xerrors.New("service is not provided")
	}

	if err != nil {
		return nil, xerrors.Errorf("provide service error: %w", err)
	}

	return serviceInstance, nil
}

// Start ハンドラを追加して監視を開始
func (s *Service) Start() error {
	s.session.AddHandler(newHandler(s))

	if err := s.session.Open(); err != nil {
		return xerrors.Errorf("session open error: %w", err)
	}

	return nil
}

// Send メッセージ送信
func (s *Service) Send(msg string) {
	if _, err := s.session.ChannelMessageSend(s.channelID, msg); err != nil {
		logger.Error(
			fmt.Sprintf("%+v", xerrors.Errorf("message send error: %w", err)),
		)
	}
}

// SendTo メッセージ送信（対象を取る）
func (s *Service) SendTo(userID, msg string) {
	s.Send(fmt.Sprintf("<@!%s>\n%s", userID, msg))
}

func (s *Service) getCommand(m *discordgo.MessageCreate) string {
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

func (s *Service) isCommand(m *discordgo.MessageCreate) bool {
	return s.botID != m.Author.ID &&
		m.ChannelID == s.channelID &&
		(strings.HasPrefix(m.Content, fmt.Sprintf("<@%s>", s.botID)) ||
			strings.HasPrefix(m.Content, fmt.Sprintf("<@!%s>", s.botID))) &&
		s.getCommand(m) != ""
}
