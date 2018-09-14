package main

import (
	"strconv"
	"math"
	"regexp"
	"github.com/bwmarrin/discordgo"
)

func clearMessages(session *discordgo.Session, m *discordgo.MessageCreate) {
	userPermissions, _ := session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	userPermissionsB := strconv.FormatInt(int64(userPermissions), 2)
	if (!(string([]rune(userPermissionsB)[15:16]) == "1")) {
		return
	}
	content := m.Content
	regex := regexp.MustCompile("\\d+")
	if (!regex.MatchString(content)) { return }
	match := regex.FindString(content)
	matchInt, _ := strconv.ParseFloat(match, 64)
	clearCount := math.Min(matchInt, 100)
	messages, _ := session.ChannelMessages(m.ChannelID, int(clearCount), "", "", "")
	var messageIDs []string
	l := len(messages)
	for i := 0; i < l; i++ {
		messageIDs[i] = messages[i].ID
	}
	session.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
}