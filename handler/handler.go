package handler

import (
	"fmt"

	"github.com/biohuns/discord-servertool/entity"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

// NewHandler ハンドラ生成
func NewHandler(
	is entity.InstanceService,
	ms entity.MessageService,
) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if !ms.IsCommand(m) {
			return
		}

		cmd := ms.GetCommand(m)

		msg, err := execute(is, cmd)
		if err != nil {
			ms.SendTo(m.Author.ID, fmt.Sprintf("```ERROR: %s``````%+v```", cmd, err))
			logger.Error(fmt.Sprintf("%+v", err))
		}

		ms.SendTo(m.Author.ID, msg)
	}
}

func execute(is entity.InstanceService, cmd string) (string, error) {
	switch cmd {
	// インスタンス起動
	case "start":
		if err := is.Start(); err != nil {
			return "", xerrors.Errorf("start error: %w", err)
		}
		return "```起動中...```", nil

	// インスタンス停止
	case "stop":
		if err := is.Stop(); err != nil {
			return "", xerrors.Errorf("stop error: %w", err)
		}
		return "```停止中...```", nil

	// インスタンス状態確認
	case "status":
		s, err := is.Status()
		if err != nil {
			return "", xerrors.Errorf("get status error: %w", err)
		}
		return fmt.Sprintf("```%s```", s.GetStatus()), nil

	case "help":
		return "", nil

	default:
		return "```不正なコマンド``````\nstart:  起動\nstop:   停止\nstatus: 状態確認```", nil
	}
}
