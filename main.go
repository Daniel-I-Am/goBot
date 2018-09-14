package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

const prefix string = "!"

var voiceSession *discordgo.VoiceConnection
var Token string
var BotID string
var err error

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	checkError(err)
	user, err := dg.User("@me")
	checkError(err)
	BotID = user.ID
	fmt.Println("Bot is running as '" + user.Username + "'")

	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	// Set up some last things
	rand.Seed(time.Now().UnixNano()) // seed RNG with time 

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, prefix) {
		fmt.Println("Command run by '" + m.Author.Username + "':", m.Content)
	} else {
		return //not a command
	}
	defer sendError(s, m.ChannelID)
	if strings.HasPrefix(m.Content, prefix + "roll") {
		s.ChannelMessageSend(m.ChannelID, roll(m.Content))
	}
	if strings.HasPrefix(m.Content, prefix + "join") {
		voiceSession, err = joinUserVoiceChannel(s, m.Author.ID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		} else {
			fmt.Println("out", voiceSession)
		}
	}
	if strings.HasPrefix(m.Content, prefix + "leave") {
		leaveVoiceChannel()
	}
	if strings.HasPrefix(m.Content, prefix + "play") {
		playVideo(s, m)
	}
	if strings.HasPrefix(m.Content, prefix + "clear") {
		clearMessages(s, m)
	}
}