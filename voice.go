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
	return nil, errors.New("I cannot see your voice channel, please join a voice channel I can see")
}

func leaveVoiceChannel(voiceSession *discordgo.VoiceConnection) {
	fmt.Println("Leaving session")
	voiceSession.Disconnect()
}

func playVideo(session *discordgo.Session, content string) {
	return
}