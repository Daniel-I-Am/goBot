package main

import (
	"fmt"
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
		fmt.Println("Clearing",n,"to",min(clearCount,n+100))
		countLeft := clearCount-n
		fmt.Println("That is",countLeft,"messages")
		messages, _ := session.ChannelMessages(m.ChannelID, int(min(countLeft, 100)), "", "", "")
		l := len(messages)
		messageIDs := make([]string, l)
		for i := 0; i < l; i++ {
			messageIDs[i] = messages[i].ID
		}
		session.ChannelMessagesBulkDelete(m.ChannelID, messageIDs)
	}
}