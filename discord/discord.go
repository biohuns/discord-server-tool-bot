package discord

import (
	"fmt"
	"strings"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/gcp"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

const (
	start   = "```STARTING...```"
	stop    = "```STOPPING...```"
	help    = "```\nHELP\n\tstart: サーバーを起動する\n\tstop: サーバーを停止する\n\tstatus : サーバーの状態を確認する```"
	unknown = "```UNKNOWN COMMAND```"
)

var (
	prefix1 string
	prefix2 string
)

// Start セッションを開く
func Start() error {
	session, err := discordgo.New()
	if err != nil {
		return xerrors.Errorf("failed to start: %w", err)
	}

	session.Token = fmt.Sprintf("Bot %s", config.Get().Discord.Token)
	session.AddHandler(onMessageCreate)

	if err := session.Open(); err != nil {
		return xerrors.Errorf("failed to start: %w", err)
	}

	prefix1 = fmt.Sprintf("<@%s>", config.Get().Discord.BotID)
	prefix2 = fmt.Sprintf("<@!%s>", config.Get().Discord.BotID)

	return nil
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot 自身のメッセージはスキップ
	if m.Author.ID == config.Get().Discord.BotID {
		return
	}

	// 対象チャンネル以外のメッセージはスキップ
	if m.ChannelID != config.Get().Discord.ChannelID {
		return
	}

	// Bot 宛以外のメッセージはスキップ
	cmd := getCommand(m.Content)
	if cmd == "" {
		return
	}

	message, err := executeCommand(cmd)
	if err != nil {
		logger.Error(fmt.Sprintf("%+v", err))
		send(m, s, fmt.Sprintf("```%s ERROR```", strings.ToUpper(cmd)))
	}

	send(m, s, message)
}

func executeCommand(cmd string) (string, error) {
	switch cmd {
	// インスタンス起動
	case "start":
		if err := gcp.Shared.Start(); err != nil {
			return "", xerrors.Errorf("failed to execute %s: %w", cmd, err)
		}
		return start, nil

	// インスタンス停止
	case "stop":
		if err := gcp.Shared.Stop(); err != nil {
			return "", xerrors.Errorf("failed to execute %s: %w", cmd, err)
		}
		return stop, nil

	// インスタンス状況取得
	case "status":
		instances, err := gcp.Shared.Instances()
		if err != nil {
			return "", xerrors.Errorf("failed to execute %s: %w", cmd, err)
		}
		text := "```STATUS"
		for _, instance := range instances {
			text += fmt.Sprintf("\n%s: %s", instance.Name, instance.Status)
		}
		text += "```"
		return text, nil

	// ヘルプ
	case "help":
		return help, nil

	// 不正なコマンド
	default:
		return unknown + help, nil
	}
}

func getCommand(message string) string {
	message = strings.TrimSpace(message)

	if !(strings.HasPrefix(message, prefix1) || strings.HasPrefix(message, prefix2)) {
		return ""
	}

	message = strings.Replace(message, prefix1, "", 1)
	message = strings.Replace(message, prefix2, "", 1)
	message = strings.TrimSpace(message)

	return message
}

func send(m *discordgo.MessageCreate, s *discordgo.Session, msg string) {
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		logger.Error(fmt.Sprintf("%+v", xerrors.Errorf("failed to send message: %w", err)))
	}
}
