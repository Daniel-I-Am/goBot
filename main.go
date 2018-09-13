package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"unicode"
	"strconv"
	"regexp"
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

const prefix string = "!"
var token string
var BotID string

func init() {
	flag.StringVar(&token, "t", "", "Bot token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + token)
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
	// If the message is "ping" reply with "Pong!"
	if strings.HasPrefix(m.Content, prefix + "roll") {
		s.ChannelMessageSend(m.ChannelID, roll(m.Content))
	}
}

func stripWhitespace(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
		  return -1
		}
		return r
	  }, str)
}

func roll(content string) string {
	input := stripWhitespace(content[len(prefix)+4:])
	regex := regexp.MustCompile("(?P<pre>.*?)(?P<roll>\\d+d\\d+)(?P<suf>.*)")
	for regex.MatchString(input) == true {
		match := regex.FindStringSubmatch(input)
		input = match[1] + parseRoll(match[2]) + match[3]
	}
	// now that we have the dice rolled, we need to calculate the result
	return input
}

func parseRoll(input string) string {
	regex := regexp.MustCompile("(?P<a>\\d+)d(?P<b>\\d+)")
	if !regex.MatchString(input) { return "" }
	match := regex.FindStringSubmatch(input)
	toRet := "( "
	diceCount, _ := strconv.Atoi(match[1])
	diceSize, _ := strconv.Atoi(match[2])
	for i := 0; i < diceCount; i++ {
		rn := rand.Intn(diceSize)
		toRet += strconv.Itoa(rn) + " + "
	}
	toRet = toRet[:len(toRet)-3]
	toRet += " )"
	return toRet
}