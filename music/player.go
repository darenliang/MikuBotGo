package music

import (
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"log"
	"runtime"
)

// MusicConnections maps a Guild ID to an associated voice connection.
var (
	MusicConnections = map[string]*Connection{}
	UserJoinError    = errors.New("you are not in a voice channel")
)

var EncOpts = &dca.EncodeOptions{
	Volume:           256,
	Channels:         2,
	FrameRate:        48000,
	FrameDuration:    20,
	Bitrate:          96,
	PacketLoss:       1,
	RawOutput:        true,
	Application:      dca.AudioApplicationLowDelay,
	CoverFormat:      "jpeg",
	CompressionLevel: 10,
	BufferedFrames:   100,
	VBR:              true,
}

func AddToQueue(ctx *exrouter.Context, conn *Connection, song *SongResponse) {
	conn.AddYouTubeVideo(song)
	ctx.Reply(fmt.Sprintf(":white_check_mark: Added %s to the queue.", song.Title))
}

func PlaySong(ctx *exrouter.Context, guildID, musicChannelID string, song *SongResponse) {
	voice, err := ctx.Ses.ChannelVoiceJoin(guildID, musicChannelID, false, false)
	if err != nil {
		log.Print("music: voice join fail")
		ctx.Reply(":warning: Failed to join voice channel.")
		return
	}

	voice.LogLevel = discordgo.LogWarning

	conn := NewConnection(voice, EncOpts)
	MusicConnections[guildID] = conn

	conn.AddYouTubeVideo(song)

	for voice.Ready == false {
		runtime.Gosched()
	}

	ctx.Reply(fmt.Sprintf(":arrow_forward: Started playing %s", song.Title))
	err = conn.StreamMusic()

	if err != nil {
		log.Printf("music: start music fail: %s", song.Title)
		return
	}

	conn.Close()
}

func GetPlayReadyData(musicChannelID string, currChannelGuildID string) (bool, *Connection, error) {
	conn, ok := MusicConnections[currChannelGuildID]
	if ok {
		if musicChannelID != conn.ChannelID {
			return false, nil, UserJoinError
		}

		if conn.Stream == nil {
			return true, conn, nil
		}

		fin, _ := conn.Stream.Finished()

		if fin {
			return true, conn, nil
		}

		return false, conn, nil
	} else {
		if musicChannelID == "" {
			return false, nil, UserJoinError
		}

		return true, conn, nil
	}
}
