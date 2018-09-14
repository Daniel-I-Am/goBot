package main

import (
	"fmt"
	"errors"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

func joinUserVoiceChannel(session *discordgo.Session, userID string) (*discordgo.VoiceConnection, error) {
	// Find a user's current voice channel
	vs, err := findUserVoiceState(session, userID)
	if err != nil {
		return nil, err
	}

	// Join the user's channel and start unmuted and deafened.
	ses, err := session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
	return ses, err
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

func leaveVoiceChannel() {
	fmt.Println("Leaving voice session")
	voiceSession.Disconnect()
}

func playVideo(session *discordgo.Session, m *discordgo.MessageCreate) {
	print("Checking URL")
	videoURL := m.Content[len(prefix)+5:]
	print("Found '" + videoURL + "'")
	if voiceSession != nil {
		print("Voice Session not ready")
		voiceSession, err = joinUserVoiceChannel(session, m.Author.ID)
		if err != nil {
			print("join channel failed")
			return
		}
		print("join channel succeeded")
	}
	print("Is in voice channel")
	// Change these accordingly
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	print("Set option up")

	videoInfo, err := ytdl.GetVideoInfo(videoURL)
	if err != nil {
		fmt.Println("Error in 1")
	}
	print("Got vid info")

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		fmt.Println("Error in 2")
	}
	print("Got download URL")

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		fmt.Println("Error in 3")
	}
	print("Encoded file")
	defer encodingSession.Cleanup()
	print("Set up cleanup")
		
	done := make(chan error)    
	print("done")
	dca.NewStream(encodingSession, voiceSession, done)
	print("new stream")
	someErr := <- done
	print("someErr")
	if someErr != nil && someErr != io.EOF {
		fmt.Println("Error in 4")
	}
	print("Finished")
}

func print(m string) {
	fmt.Println(m)
}