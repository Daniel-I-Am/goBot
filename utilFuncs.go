package main

import (
	"strings"
	"fmt"
	"time"
	"strconv"
	"github.com/bwmarrin/discordgo"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func sendError(ses *discordgo.Session, cid string) {
	//log(err)
	if err := recover(); err != nil {
		ses.ChannelMessageSend(cid, "Sorry, something went wrong.\nhttps://media.makeameme.org/created/you-caused-a-5b9ab5.jpg")
	}
}

func noPermission(session *discordgo.Session, m *discordgo.MessageCreate, command string) {
	message, _ := session.ChannelMessageSend(m.ChannelID, "You do not have the permission to use the `" + command + "` command")
	session.ChannelTyping(m.ChannelID)
	time.Sleep(5 * time.Second)
	session.ChannelMessageDelete(m.ChannelID, message.ID)
	session.ChannelMessageDelete(m.ChannelID, m.ID)
}

func checkPermission(session *discordgo.Session, m *discordgo.MessageCreate, n int, command string) bool {
	userPermissions, _ := session.UserChannelPermissions(m.Author.ID, m.ChannelID)
	userPermissionsB := strconv.FormatInt(int64(userPermissions), 2)
	if (!(string([]rune(userPermissionsB)[n:n+1]) == "1")) {
		// missing permission
		noPermission(session, m, command)
		return false
	}
	return true
}

func log(m ...string) {
	SendToSTDOUT("[stdout]",strings.Join(m," "))
}

func logErr(m ...string) {
	SendToSTDOUT("[stderr]",strings.Join(m," "))
}

func SendToSTDOUT(m ...string) {
	fmt.Println("[GoBot]",strings.Join(m, " "))
}