package handler

import (
	"encoding/json"
	"fmt"

	"github.com/biohuns/discord-servertool/adapter/discord"
	"github.com/biohuns/discord-servertool/adapter/gcp"
	"github.com/biohuns/discord-servertool/logger"
	"github.com/bwmarrin/discordgo"
	"golang.org/x/xerrors"
)

// NewHandler ハンドラ生成
func NewHandler(gs gcp.Service, ds discord.Service) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		b, _ := json.Marshal(m)
		logger.Debug(string(b))
		fmt.Println(ds.IsCommand(m))

		if !ds.IsCommand(m) {
			return
		}

		cmd := ds.GetCommand(m)

		msg, err := execute(gs, cmd)
		if err != nil {
			ds.Send(fmt.Sprintf("```ERROR: %s``````%+v```", cmd, err))
			logger.Error(fmt.Sprintf("%+v", err))
		}

		ds.Send(msg)
	}
}

func execute(gs gcp.Service, cmd string) (string, error) {
	switch cmd {
	// インスタンス起動
	case "start":
		if err := gs.Start(); err != nil {
			return "", xerrors.Errorf("start error: %w", err)
		}
		return "```起動中...```", nil

	// インスタンス停止
	case "stop":
		if err := gs.Stop(); err != nil {
			return "", xerrors.Errorf("stop error: %w", err)
		}
		return "```停止中...```", nil

	// インスタンス状態確認
	case "status":
		s, err := gs.Status()
		if err != nil {
			return "", xerrors.Errorf("get status error: %w", err)
		}
		return fmt.Sprintf("```%s: %s```", s.Name, s.Status), nil

	case "help":
		return "", nil

	default:
		return "```不正なコマンド``````\nHELP\n\tstart: 起動\n\tstop: 停止\n\tstatus : 状態確認```", nil
	}
}
