package main

import (
	"flag"
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
var session *discordgo.Session
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
	log("Bot is running as '" + user.Username + "'")

	dg.AddHandler(messageCreate)
	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log("error opening connection,", err.Error())
		return
	}
	// Set up some last things
	rand.Seed(time.Now().UnixNano()) // seed RNG with time 

	// Wait here until CTRL-C or other term signal is received.
	log("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(ses *discordgo.Session, m *discordgo.MessageCreate) {
	session = ses
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == session.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, prefix) {
		log("Command run by '" + m.Author.Username + "':", m.Content)
	} else {
		return //not a command
	}
	defer sendError(m.ChannelID)
	if strings.HasPrefix(m.Content, prefix + "roll") {
		session.ChannelMessageSend(m.ChannelID, roll(m.Content))
	}
	if strings.HasPrefix(m.Content, prefix + "join") {
		voiceSession, err = joinUserVoiceChannel(m.Author.ID)
		if err != nil {
			session.ChannelMessageSend(m.ChannelID, err.Error())
		} else {
			log("Connected to voice session")
		}
	}
	if strings.HasPrefix(m.Content, prefix + "leave") {
		leaveVoiceChannel()
	}
	if strings.HasPrefix(m.Content, prefix + "play") {
		playVideo(m)
	}
	if strings.HasPrefix(m.Content, prefix + "clear") {
		clearMessages(m)
	}
	if strings.HasPrefix(m.Content, prefix + "friend") {
		session.RelationshipFriendRequestSend(m.Author.ID)
	}
}