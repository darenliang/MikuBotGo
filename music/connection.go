package music

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
	"io"
	"sync"
)

// Connection is the type for a music connection to Discord.
type Connection struct {
	GuildID         string
	ChannelID       string
	EncodeOpts      *dca.EncodeOptions
	VoiceConnection *discordgo.VoiceConnection
	Stream          *dca.StreamingSession
	Done            chan error
	Queue           []*SongResponse
	Mutex           *sync.Mutex
}

// NewConnection will return a new Connection struct.
func NewConnection(voice *discordgo.VoiceConnection, opts *dca.EncodeOptions) *Connection {
	return &Connection{
		GuildID:         voice.GuildID,
		ChannelID:       voice.ChannelID,
		EncodeOpts:      opts,
		VoiceConnection: voice,
		Mutex:           &sync.Mutex{},
	}
}

// StreamMusic will create a new encode session from the current DownloadURL
// and stream that to the VoiceConnection.
// Will block until queue is empty.
func (c *Connection) StreamMusic() error {
	for len(c.Queue) != 0 {
		c.Mutex.Lock()
		encodeSession, err := dca.EncodeFile(c.Queue[0].URL, c.EncodeOpts)
		c.Mutex.Unlock()
		if err != nil {
			return err
		}
		c.Mutex.Lock()
		c.Done = make(chan error)
		c.Stream = dca.NewStream(encodeSession, c.VoiceConnection, c.Done)
		c.Mutex.Unlock()
		derr := <-c.Done
		if derr != nil && derr != io.EOF {
			encodeSession.Cleanup()
			return derr
		}
		encodeSession.Cleanup()
		c.Mutex.Lock()
		c.Queue = c.Queue[1:]
		c.Mutex.Unlock()
	}

	return nil
}

// AddYouTubeVideo will add the download URL for a YouTube video to the queue.
func (c *Connection) AddYouTubeVideo(song *SongResponse) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Queue = append(c.Queue, song)
}

// Close closes the VoiceConnection, stops sending speaking packet, and closes
// the EncodeSession.
func (c *Connection) Close() error {
	err := c.VoiceConnection.Speaking(false)
	if err != nil {
		return err
	}

	c.VoiceConnection.Close()
	err = c.VoiceConnection.Disconnect()
	if err != nil {
		return err
	}

	delete(MusicConnections, c.GuildID)

	return nil
}
