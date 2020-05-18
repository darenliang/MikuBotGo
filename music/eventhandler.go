package music

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/foxbot/gavalink"
	"log"
	"os"
)

type EventHandler struct{}

func (m EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	if len(AudioPlayers[player.GuildID()].Queue) == 0 {
		_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, "Queue ended.")
		AudioPlayers[player.GuildID()].Playing = false
		return nil
	}

	var idx int
	var val gavalink.Track
	for idx, val = range AudioPlayers[player.GuildID()].Queue {
		if err := AudioPlayers[player.GuildID()].Player.Play(val.Data); err != nil {
			_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, fmt.Sprintf("Failed to play: %s", val.Info.Title))
		} else {
			_, _ = Session.ChannelMessageSendEmbed(AudioPlayers[player.GuildID()].ChannelID, &discordgo.MessageEmbed{
				Color:       config.EmbedColor,
				Title:       "Now playing",
				Description: val.Info.Title,
				URL:         val.Info.URI,
			})
			break
		}
	}

	AudioPlayers[player.GuildID()].Queue = AudioPlayers[player.GuildID()].Queue[idx+1:]
	return nil
}

func (m EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	return nil
}

func (m EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	return nil
}

type GuildPlayer struct {
	Player    *gavalink.Player
	ChannelID string
	Playing   bool
	Queue     []gavalink.Track
}

var (
	Session       *discordgo.Session
	AudioLavalink *gavalink.Lavalink
	AudioNode     *gavalink.Node
	AudioPlayers  map[string]*GuildPlayer
	lavalinkRest  string
	lavalinkWS    string
	lavalinkPass  string
)

func init() {
	lavalinkRest = os.Getenv("LAVALINK_REST")
	lavalinkWS = os.Getenv("LAVALINK_WS")
	lavalinkPass = os.Getenv("LAVALINK_PASS")
}

func AudioInit(botID string) {
	AudioLavalink = gavalink.NewLavalink("1", botID)
	AudioPlayers = make(map[string]*GuildPlayer)

	err := AudioLavalink.AddNodes(gavalink.NodeConfig{
		REST:      lavalinkRest,
		WebSocket: lavalinkWS,
		Password:  lavalinkPass,
	})

	if err != nil {
		log.Println(err)
	}
}
