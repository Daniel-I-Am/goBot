package main

import (
	"fmt"
	"time"
	"github.com/bwmarrin/discordgo"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func sendError(ses *discordgo.Session, cid string) {
	//fmt.Println(err)
	if err := recover(); err != nil {
		ses.ChannelMessageSend(cid, "Sorry, something went wrong.\nhttps://media.makeameme.org/created/you-caused-a-5b9ab5.jpg")
	}
}

func noPermission(session *discordgo.Session, m *discordgo.MessageCreate, command string) {
	message, _ := session.ChannelMessageSend(m.ChannelID, "You do not have the permission to use the `clear` command")
	session.ChannelTyping(m.ChannelID)
	time.Sleep(5 * time.Second)
	session.ChannelMessageDelete(m.ChannelID, message.ID)
	session.ChannelMessageDelete(m.ChannelID, m.ID)
}

func checkPermission() {
}