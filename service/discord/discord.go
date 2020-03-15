package discord

import (
	"fmt"
	"strings"
	"sync"

	"github.com/biohuns/discord-servertool/util"

	"github.com/biohuns/discord-servertool/entity"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

// Service Discordサービス
type Service struct {
	log       entity.LogService
	instance  entity.InstanceService
	server    entity.ServerStatusService
	session   *discordgo.Session
	channelID string
	botID     string
}

// Start ハンドラを追加して監視を開始
func (s *Service) Start() error {
	s.session.AddHandler(s.newHandler())

	if err := s.session.Open(); err != nil {
		return xerrors.Errorf("failed to open session: %w", err)
	}

	return nil
}

// Send メッセージ送信
func (s *Service) Send(userID, msg string) error {
	if userID != "" {
		msg = fmt.Sprintf("<@!%s>\n%s", userID, msg)
	}

	if _, err := s.session.ChannelMessageSend(s.channelID, msg); err != nil {
		return xerrors.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (s *Service) newHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(sess *discordgo.Session, m *discordgo.MessageCreate) {
		if !s.isCommand(m) {
			return
		}

		switch s.getCommand(m) {
		// インスタンス起動
		case "start":
			if err := s.instance.Start(); err != nil {
				_ = s.Send(m.Author.ID, fmt.Sprintf("```インスタンス起動処理失敗``````%+v```", err))
				s.log.Error(xerrors.Errorf("failed to start instance: %w", err))
			}
			_ = s.Send(m.Author.ID, "```起動処理中...```")

		// インスタンス停止
		case "stop":
			if err := s.instance.Stop(); err != nil {
				_ = s.Send(m.Author.ID, fmt.Sprintf("```インスタンス停止処理失敗``````%+v```", err))
				s.log.Error(xerrors.Errorf("failed to stop instance: %w", err))
			}
			_ = s.Send(m.Author.ID, "```停止処理中...```")

		// インスタンス状態確認
		case "status":
			instanceStatus, err := s.instance.Status()
			if err != nil {
				_ = s.Send(m.Author.ID, fmt.Sprintf("```インスタンス状態確認失敗``````%+v```", err))
				s.log.Error(xerrors.Errorf("failed to get instance: %w", err))
			}
			if instanceStatus.Status == entity.InstanceStatusRunning {
				serverStatus, err := s.server.Status()
				if err != nil {
					_ = s.Send(m.Author.ID, fmt.Sprintf("```サーバ状態確認失敗``````%+v```", err))
					s.log.Error(xerrors.Errorf("failed to get server status: %w", err))
				}

				_ = s.Send(m.Author.ID,
					util.InstanceStatusText(
						instanceStatus.Name,
						instanceStatus.Status.String(),
					)+util.ServerStatusText(
						serverStatus.IsOnline,
						serverStatus.GameName,
						serverStatus.PlayerCount,
						serverStatus.MaxPlayerCount,
						serverStatus.Map,
					))
			} else {
				_ = s.Send(m.Author.ID, util.InstanceStatusText(
					instanceStatus.Name,
					instanceStatus.Status.String(),
				))
			}

		default:
			_ = s.Send(m.Author.ID, "```start:  起動\nstop:   停止\nstatus: 状態確認```")
		}
	}
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

var (
	shared *Service
	once   sync.Once
)

// ProvideService サービス返却
func ProvideService(
	log entity.LogService,
	cache entity.ConfigService,
	instance entity.InstanceService,
	server entity.ServerStatusService,
) (entity.MessageService, error) {
	var err error

	once.Do(func() {
		var session *discordgo.Session
		session, err = discordgo.New()
		if err != nil {
			err = xerrors.Errorf("failed to create session: %w", err)
			return
		}

		token, channelID, botID := cache.GetDiscordConfig()
		session.Token = fmt.Sprintf("Bot %s", token)

		shared = &Service{
			log:       log,
			instance:  instance,
			server:    server,
			session:   session,
			channelID: channelID,
			botID:     botID,
		}
	})

	if shared == nil {
		err = xerrors.Errorf("service is not provided: %w", err)
	}

	if err != nil {
		return nil, xerrors.Errorf("failed to provide service: %w", err)
	}

	return shared, nil
}
