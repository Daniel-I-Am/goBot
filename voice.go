package main

import (
	"fmt"
	"errors"

	"github.com/bwmarrin/discordgo"
)

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
	fmt.Println("Leaving session")
	voiceSession.Close()
}

func playVideo(session *discordgo.Session, content string) {
	return
}