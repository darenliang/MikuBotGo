package music

import (
	"github.com/bwmarrin/discordgo"
	"sync"
	"time"
)

type Connection struct {
	voiceConnection *discordgo.VoiceConnection
	send            chan []int16
	lock            sync.Mutex
	sendpcm         bool
	stopRunning     bool
	playing         bool
}

func NewConnection(voiceConnection *discordgo.VoiceConnection) *Connection {
	connection := new(Connection)
	connection.voiceConnection = voiceConnection
	return connection
}
func (connection *Connection) Disconnect() {
	connection.voiceConnection.Disconnect()
}

func (connection *Connection) IsPlaying() bool {
	return connection.playing
}

// HandleMusicTimeout handles idle player and disconnects after 15 minutes of IDLE
func HandleMusicTimeout(sess *Session, callback func(string)) {
	startTime := time.Now()
	for {
		if sess.Connection.IsPlaying() {
			return
		}
		if time.Since(startTime).Minutes() > 15 {
			musicSession := MusicSessions.GetByGuild(sess.guildId)

			if musicSession == nil {
				return
			}

			// TODO: Fix the type signature
			MusicSessions.Leave(nil, *musicSession)

			callback("Disconnected because of inactivity.")
			return
		}
		time.Sleep(1)
	}
}
