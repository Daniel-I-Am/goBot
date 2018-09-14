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
	"errors"

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
	rand.Seed(time.Now().UnixNano())

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	defer sendError(s, m.ChannelID)
	// If the message is "ping" reply with "Pong!"
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
	if strings.HasPrefix(m.Content, prefix+"leave") {
		leaveVoiceChannel(voiceSession)
	}
}

func sendError(ses *discordgo.Session, cid string) {
	if err := recover(); err != nil {
		ses.ChannelMessageSend(cid, "Sorry, something went wrong.\nhttps://media.makeameme.org/created/you-caused-a-5b9ab5.jpg")
	}
}

func joinUserVoiceChannel(session *discordgo.Session, userID string) (*discordgo.VoiceConnection, error) {
	// Find a user's current voice channel
	vs, err := findUserVoiceState(session, userID)
	if err != nil {
		return nil, err
	}

	// Join the user's channel and start unmuted and deafened.
	return session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
}

func findUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}

func leaveVoiceChannel(voiceSession *discordgo.VoiceConnection) {
	voiceSession.Close()
}
