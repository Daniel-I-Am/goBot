package main

import (
	"strconv"
	"regexp"
	"github.com/bwmarrin/discordgo"
)

func min(a int, b int) int {
	if a > b {
		return b
	}
	return a
}

func clearMessages(session *discordgo.Session, m *discordgo.MessageCreate) {
	if !checkPermission(session, m, 13, "clear") { return }
	content := m.Content
	regex := regexp.MustCompile("\\d+")
	if (!regex.MatchString(content)) { return }
	match := regex.FindString(content)
	matchInt, _ := strconv.ParseFloat(match, 64)
	clearCount := int(matchInt) + 1 // 1 since we need 
	for n := 0; n<clearCount; n+=100 {
		log("Clearing message", string(n), "to", string(min(clearCount,n+100)))
		countLeft := min(clearCount-n, 100)
		log("That is", string(countLeft),"messages")
		messages, _ := session.ChannelMessages(m.ChannelID, countLeft, "", "", "")
		l := len(messages)
		messageIDs := make([]string, l)
		for i := 0; i < l; i++ {
			messageIDs[i] = messages[i].ID
		}
		session.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	}
	messageTag := "messages"
	if clearCount == 2 { messageTag = "message"}
	session.ChannelMessageSend(m.ChannelID, "Cleared " + strconv.Itoa(clearCount - 1) + " " + messageTag + " requested by <@" + m.Author.ID + ">")
}