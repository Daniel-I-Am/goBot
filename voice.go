package main

import (
	"errors"
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"github.com/rylio/ytdl"
)

var volume int

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
	log("Leaving voice session")
	voiceSession.Disconnect()
}

func playVideo(session *discordgo.Session, m *discordgo.MessageCreate) {
	log("Checking URL")
	videoURL := m.Content[len(prefix)+5:]
	log("Found '" + videoURL + "'")
	if voiceSession != nil {
		log("Voice Session not ready")
		voiceSession, err = joinUserVoiceChannel(session, m.Author.ID)
		if err != nil {
			log("join channel failed")
			return
		}
		log("join channel succeeded")
	}
	log("Is in voice channel")
	// Change these accordingly
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"
	log("Set option up")

	videoInfo, err := ytdl.GetVideoInfo(videoURL)
	if err != nil {
		log("Error in 1")
	}
	log("Got vid info")

	format := videoInfo.Formats.Extremes(ytdl.FormatAudioBitrateKey, true)[0]
	downloadURL, err := videoInfo.GetDownloadURL(format)
	if err != nil {
		log("Error in 2")
	}
	log("Got download URL")

	encodingSession, err := dca.EncodeFile(downloadURL.String(), options)
	if err != nil {
		log("Error in 3")
	}
	log("Encoded file")
	defer encodingSession.Cleanup()
	log("Set up cleanup")
		
	done := make(chan error)    
	log("done")
	dca.NewStream(encodingSession, voiceSession, done)
	log("new stream")
	someErr := <- done
	log("someErr")
	if someErr != nil && someErr != io.EOF {
		log("Error in 4")
	}
	log("Finished")
}