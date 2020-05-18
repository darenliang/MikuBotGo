package music

import (
	"bytes"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/foxbot/gavalink"
	"github.com/robfig/cron"
	"log"
	"os"
)

const stuckTimeout = 10 * 1000

type EventHandler struct{}

func (m EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player) {
		return nil
	}

	PlayNextTrack(player)
	return nil
}

func (m EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player) {
		return nil
	}

	_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, fmt.Sprintf("Error on playing: %s", AudioPlayers[player.GuildID()].CurrTrack.Info.Title))
	PlayNextTrack(player)
	return nil
}

func (m EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	if IsEmptyQueue(player) {
		return nil
	}

	if threshold > stuckTimeout {
		_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, fmt.Sprintf("Stuck on playing: %s", AudioPlayers[player.GuildID()].CurrTrack.Info.Title))
		PlayNextTrack(player)
	}
	return nil
}

func PlayNextTrack(player *gavalink.Player) {
	var idx int
	var val gavalink.Track
	AudioPlayers[player.GuildID()].Playing = false
	for idx, val = range AudioPlayers[player.GuildID()].Queue {
		if err := AudioPlayers[player.GuildID()].Player.Play(val.Data); err != nil {
			_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, fmt.Sprintf("Failed to play: %s", val.Info.Title))
		} else {
			AudioPlayers[player.GuildID()].Playing = true
			AudioPlayers[player.GuildID()].CurrTrack = val
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

	if !AudioPlayers[player.GuildID()].Playing {
		_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, "Queue ended.")
	}
}

func IsEmptyQueue(player *gavalink.Player) bool {
	// Set to not playing initially
	AudioPlayers[player.GuildID()].Playing = false

	if len(AudioPlayers[player.GuildID()].Queue) == 0 {
		_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, "Queue ended.")
		AudioPlayers[player.GuildID()].Playing = false
		return true
	}
	return false
}

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
	var (
		buf      = bytes.NewBuffer([]byte{})
		c        = cron.New()
		prevcont = buf.String()
	)

	gavalink.Log = log.New(buf, "[gavalink] ", 0)

	if err := c.AddFunc("@every 500ms", func() {
		if buf.String() == prevcont {
			return
		}

		log.Printf("%s", buf.String())
		buf.Reset()
	}); err != nil {
		log.Printf("Initializing a logger for gavalink failed: \"%s\".", err.Error())
	}

	c.Start()

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
