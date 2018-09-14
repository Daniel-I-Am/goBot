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
	clearCount := int(matchInt)
	for n := 0; n<clearCount; n+=100 {
		countLeft := float64(clearCount-n)
		messages, _ := session.ChannelMessages(m.ChannelID, int(math.Min(countLeft, 100)), "", "", "")
		l := len(messages)
		messageIDs := make([]string, l)
		for i := 0; i < l; i++ {
			messageIDs[i] = messages[i].ID
		}
		session.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	}
}