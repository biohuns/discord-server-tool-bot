package discord

import (
	"fmt"
	"log"
	"strings"

	"github.com/biohuns/discord-servertool/config"
	"github.com/biohuns/discord-servertool/gcp"
	"github.com/bwmarrin/discordgo"
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
		return
	}
	if m.ChannelID != config.Get().Discord.ChannelID {
		log.Println("other channel message")
		return
	}

	if !strings.HasPrefix(m.Content, command("")) {
		log.Println("non-command message")
		return
	}

	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("error getting channel: ", err)
		return
	}

	log.Printf("%s > %s\n", m.Author.Username,
		strings.ReplaceAll(
			m.Content,
			command(""),
			fmt.Sprintf("@%s ", config.Get().Discord.BotName),
		),
	)

	switch {
	case strings.HasPrefix(m.Content, command("start")):
		if err := gcp.StartInstance(); err != nil {
			sendMessage(s, c, fmt.Sprintf("```\nSTART ERROR:\n%s\n```", err))
		}
		sendMessage(s, c, "```\nSTARTING...\n```")
	case strings.HasPrefix(m.Content, command("stop")):
		if err := gcp.StopInstance(); err != nil {
			sendMessage(s, c, fmt.Sprintf("```\nSTOP ERROR:\n%s\n```", err))
		}
		sendMessage(s, c, "```\nSTOPPING...\n```")
	case strings.HasPrefix(m.Content, command("status")):
		list, err := gcp.ListInstances()
		if err != nil {
			sendMessage(s, c, fmt.Sprintf("```\nSTATUS ERROR:\n%s\n```", err))
		}
		text := "```\nSTATUS"
		for _, instance := range list.Items {
			text += fmt.Sprintf("\n%s: %s", instance.Name, instance.Status)
		}
		text += "\n```"
		sendMessage(s, c, text)
	default:
		sendMessage(s, c, "```\nHELP\n\tstart: サーバーを起動する\n\tstop: サーバーを停止する\n\tstatus: サーバーの状態を確認する\n```")
	}
}

func command(name string) string {
	return fmt.Sprintf("<@!%s> %s", config.Get().Discord.BotID, name)
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("error sending message: ", err)
	}
}
