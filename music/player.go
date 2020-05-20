package music

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"log"
	"runtime"
)

// MusicConnections maps a Guild ID to an associated voice connection.
var MusicConnections = map[string]*Connection{}

var EncOpts = &dca.EncodeOptions{
	Volume:           256,
	Channels:         2,
	FrameRate:        48000,
	FrameDuration:    20,
	Bitrate:          64,
	PacketLoss:       1,
	RawOutput:        true,
	Application:      dca.AudioApplicationAudio,
	CoverFormat:      "jpeg",
	CompressionLevel: 0,
	BufferedFrames:   100,
	VBR:              true,
	Threads:          0,
	AudioFilter:      "",
	Comment:          "",
}

func AddToQueue(ctx *exrouter.Context, conn *Connection, song string) {
	vid, err := conn.AddYouTubeVideo(song)
	if err != nil {
		log.Printf("music: add song fail: %s", song)
		ctx.Reply("The requested song(s) are not available.")
		return
	}

	ctx.Reply(fmt.Sprintf("Added %s to the queue.", vid.Title))
}

func PlaySong(ctx *exrouter.Context, guildID, musicChannelID, song string) {
	voice, err := ctx.Ses.ChannelVoiceJoin(guildID, musicChannelID, false, true)
	if err != nil {
		log.Print("music: voice join fail")
		ctx.Reply("Failed to join voice channel.")
		return
	}

	voice.LogLevel = discordgo.LogWarning

	conn := NewConnection(voice, EncOpts)
	MusicConnections[guildID] = conn

	vid, err := conn.AddYouTubeVideo(song)
	if err != nil {
		log.Printf("music: add song fail: %s", song)
		ctx.Reply("The requested song(s) are not available.")
		return
	}

	for voice.Ready == false {
		runtime.Gosched()
	}

	ctx.Reply(fmt.Sprintf("Started playing %s", vid.Title))
	err = conn.StreamMusic()
	if err != nil {
		log.Printf("music: start music fail: %s", vid.Title)
		ctx.Reply("Failed to start music.")
		return
	}

	conn.Close()
}
