package music

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/foxbot/gavalink"
	"log"
	"os"
)

type GuildPlayer struct {
	Player    *gavalink.Player
	ChannelID string
	Playing   bool
	CurrTrack gavalink.Track
	Queue     []gavalink.Track
}

var (
	Session       *discordgo.Session
	AudioLavalink *gavalink.Lavalink
	AudioNode     *gavalink.Node
	AudioPlayers  map[string]*GuildPlayer
	LavalinkRest  string
	LavalinkWS    string
	LavalinkPass  string
)

func init() {
	LavalinkRest = os.Getenv("LAVALINK_REST")
	LavalinkWS = os.Getenv("LAVALINK_WS")
	LavalinkPass = os.Getenv("LAVALINK_PASS")
}

func AudioInit(botID string) {
	AudioLavalink = gavalink.NewLavalink("1", botID)
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

func Disconnect(ctx *exrouter.Context) {
	AudioPlayers[ctx.Msg.GuildID].Player.Destroy()
	ctx.Ses.ChannelVoiceJoinManual(ctx.Msg.GuildID, "", false, false)
}
