package discord

import (
	"fmt"
	"strings"

	"github.com/biohuns/discord-servertool/logger"
	"google.golang.org/api/compute/v1"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/gcp"
	"github.com/bwmarrin/discordgo"
)

const (
	start = "```STARTING...```"
	stop  = "```STOPPING...```"
	help  = "```" + `
HELP
	start: サーバーを起動する
	stop: サーバーを停止する
	status : サーバーの状態を確認する
` + "```"
	unknown = "```UNKNOWN COMMAND```"
)

func Init() error {
	session, err := discordgo.New()
	if err != nil {
		return err
	}

	session.Token = fmt.Sprintf("Bot %s", config.Get().Discord.Token)
	session.AddHandler(onMessageCreate)

	if err := session.Open(); err != nil {
		return err
	}

	return nil
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == config.Get().Discord.BotID {
		logger.Debug(logMessage(m, "bot's own message"))
		return
	}

	if m.ChannelID != config.Get().Discord.ChannelID {
		logger.Debug(logMessage(m, "invalid channel: "+m.ChannelID))
		return
	}

	prefix := fmt.Sprintf("<@!%s>", config.Get().Discord.BotID)
	command := strings.TrimSpace(strings.Replace(m.Content, prefix, "", 1))
	if !strings.HasPrefix(m.Content, prefix) {
		logger.Warn(logMessage(m, "not command: "+command))
		return
	}

	switch command {
	case "start":
		if err := gcp.StartInstance(); err != nil {
			sendError(m, s, command, err)
			return
		}
		send(m, s, start)
	case "stop":
		if err := gcp.StopInstance(); err != nil {
			sendError(m, s, command, err)
			return
		}
		send(m, s, stop)
	case "status":
		instances, err := gcp.ListInstances()
		if err != nil {
			sendError(m, s, command, err)
			return
		}
		send(m, s, status(instances.Items))
	case "help":
		send(m, s, help)
	default:
		send(m, s, unknown+help)
	}
}

func status(items []*compute.Instance) string {
	text := "```STATUS"
	for _, instance := range items {
		text += fmt.Sprintf("\n%s: %s", instance.Name, instance.Status)
	}
	text += "```"
	return text
}

func send(m *discordgo.MessageCreate, s *discordgo.Session, msg string) {
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		logMessage(m, "failed to send message")
	}
}

func sendError(
	m *discordgo.MessageCreate,
	s *discordgo.Session,
	com string,
	err error,
) {
	logger.Error(
		logMessage(m, fmt.Sprintf("failed to execute %s: %s", com, err)),
	)
	send(m, s, fmt.Sprintf("```%s ERROR```", strings.ToUpper(com)))
}

func logMessage(m *discordgo.MessageCreate, message string) string {
	return fmt.Sprintf(
		"(%s) %s >>> %s",
		message,
		m.Author.Username,
		strings.ReplaceAll(m.Content, "\n", "\\n"),
	)
}
