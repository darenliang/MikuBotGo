package music

import (
	"container/list"
	"errors"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"log"
	"sync"
	"time"
)

type Connection struct {
	*sync.RWMutex
	VoiceChannelID   string
	MsgChannelID     string
	Playing          bool
	Done             chan error
	Queue            *list.List
	VoiceConnection  *discordgo.VoiceConnection
	StreamingSession *dca.StreamingSession
	LastInvoke       time.Time
}

// AddYouTubeVideo will add the download URL for a YouTube video to the queue.
func (c *Connection) AddYouTubeVideo(song *SongResponse) {
	c.Lock()
	defer c.Unlock()
	c.Queue.PushBack(song)
}

func CreateVoiceConnection(ctx *exrouter.Context, guild *discordgo.Guild, conn *Connection) (*discordgo.VoiceConnection, error) {
	for _, vs := range guild.VoiceStates {
		if vs.UserID == ctx.Msg.Author.ID && (vs.ChannelID == conn.VoiceChannelID || !conn.Playing) {
			vc, err := ctx.Ses.ChannelVoiceJoin(guild.ID, vs.ChannelID, false, false)
			if err != nil {
				log.Print("connection: failed to join voice channel")
				ctx.Reply(":warning: Failed to join voice channel.")
				return nil, err
			}
			return vc, nil
		}
	}
	ctx.Reply(":information_source: Please join a voice channel.")
	return nil, errors.New("not in voice channel")
}

func (c *Connection) YoutubeCleanup() {
	c.Lock()
	defer c.Unlock()
	c.VoiceConnection.Disconnect()
	c.Queue = list.New()
	c.Done = make(chan error)
}
