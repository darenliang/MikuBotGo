package music

import (
	"container/list"
	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
	"log"
	"os"
)

type GuildPlayer struct {
	ChannelID string
	Queue     *list.List
}

var (
	Session       *discordgo.Session
	AudioLavalink *gavalink.Lavalink
	// TODO make thread safe
	AudioPlayers map[string]*GuildPlayer
	LavalinkRest string
	LavalinkWS   string
	LavalinkPass string
)

func init() {
	LavalinkRest = os.Getenv("LAVALINK_REST")
	LavalinkWS = os.Getenv("LAVALINK_WS")
	LavalinkPass = os.Getenv("LAVALINK_PASS")
}

func AudioInit(ses *discordgo.Session) {
	Session = ses
	AudioLavalink = gavalink.NewLavalink("1", ses.State.User.ID)
	AudioPlayers = make(map[string]*GuildPlayer)

	err := AudioLavalink.AddNodes(gavalink.NodeConfig{
		REST:      LavalinkRest,
		WebSocket: LavalinkWS,
		Password:  LavalinkPass,
	})

	if err != nil {
		log.Println(err)
	}
}

func Disconnect(guildID string) {
	delete(AudioPlayers, guildID)
	Session.ChannelVoiceJoinManual(guildID, "", false, false)
	guildPlayer, err := AudioLavalink.GetPlayer(guildID)
	if err != nil {
		return
	}
	guildPlayer.Destroy()
}
