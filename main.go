package main

import (
	"discord-ito/game"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Config struct {
	Prefix string
	DiscordToken string
}

var conf Config
var sessions map[string]*game.Session

func init() {
	sessions = make(map[string]*game.Session)
	c, _ := ini.Load("config.ini")
	conf = Config{
		Prefix: c.Section("config").Key("prefix").MustString("$"),
		DiscordToken: c.Section("config").Key("token").String(),
	}
}

func main() {
	discord, err := discordgo.New("Bot " + conf.DiscordToken)
	if err != nil {
		log.Fatal(err)
	}
	discord.AddHandler(func(s *discordgo.Session, mc *discordgo.MessageCreate) {
		if mc.Author.ID == s.State.User.ID || !strings.HasPrefix(mc.Content, conf.Prefix) {
			return
		}
		command := strings.TrimPrefix(mc.Content, conf.Prefix)
		switch command {
		case "start":
			_, ok := sessions[mc.ChannelID]
			if ok {
				_, err := s.ChannelMessageSend(mc.ChannelID, "ゲームは既に開始されています")
				if err != nil {
					log.Println(err)
				}
				return
			}
			sessions[mc.ChannelID] = game.NewSession()
			_, err := s.ChannelMessageSend(mc.ChannelID, "ゲームを開始します")
			if err != nil {
				log.Println(err)
			}
		case "join":
			session, ok := sessions[mc.ChannelID]
			if !ok {
				_, err := s.ChannelMessageSend(mc.ChannelID, conf.Prefix + "start でゲームを開始してください")
				if err != nil {
					log.Println(err)
				}
				return
			}
			num := session.JoinUser(mc.Author)
			dmChannel, err := s.UserChannelCreate(mc.Author.ID)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = s.ChannelMessageSend(dmChannel.ID, fmt.Sprintf("あなたの数字: %d", num))
			if err != nil {
				log.Println(err)
			}
		case "open":
			session, ok := sessions[mc.ChannelID]
			if !ok {
				_, err := s.ChannelMessageSend(mc.ChannelID, conf.Prefix + "start でゲームを開始してください")
				if err != nil {
					log.Println(err)
				}
				return
			}
			num, err := session.GetPlayerNumber(mc.Author.ID)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = s.ChannelMessageSend(mc.ChannelID, fmt.Sprintf("%sの数字: %d", mc.Author.Username, num))
			if err != nil {
				log.Println(err)
			}
		case "end":
			_, ok := sessions[mc.ChannelID]
			if !ok {
				_, err := s.ChannelMessageSend(mc.ChannelID, conf.Prefix + "start でゲームを開始してください")
				if err != nil {
					log.Println(err)
				}
				return
			}
			delete(sessions, mc.ChannelID)
			_, err := s.ChannelMessageSend(mc.ChannelID, "ゲームを終了しました")
			if err != nil {
				log.Println(err)
			}
		}
	})
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	err = discord.Close()
	if err != nil {
		log.Fatal(err)
	}
}