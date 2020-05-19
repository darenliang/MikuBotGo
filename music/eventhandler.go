package music

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/foxbot/gavalink"
)

const stuckTimeout = 10 * 1000

type EventHandler struct{}

func (m EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player, reason) {
		return nil
	}

	if reason != "STOP" {
		PlayNextTrack(player)
	}
	return nil
}

func (m EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player, reason) {
		return nil
	}

	_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, fmt.Sprintf("Error on playing: %s", AudioPlayers[player.GuildID()].CurrTrack.Info.Title))
	PlayNextTrack(player)
	return nil
}

func (m EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	if IsEmptyQueue(player, "") {
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
				Title:       "Playing",
				Description: val.Info.Title,
				URL:         val.Info.URI,
			})
			break
		}
	}

	if len(AudioPlayers[player.GuildID()].Queue) != 0 {
		AudioPlayers[player.GuildID()].Queue = AudioPlayers[player.GuildID()].Queue[idx+1:]
	}

	if !AudioPlayers[player.GuildID()].Playing {
		_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, "Queue ended.")
	}
}

func IsEmptyQueue(player *gavalink.Player, reason string) bool {
	// Set to not playing initially
	AudioPlayers[player.GuildID()].Playing = false

	if len(AudioPlayers[player.GuildID()].Queue) == 0 {
		if reason != "STOPPED" {
			_, _ = Session.ChannelMessageSend(AudioPlayers[player.GuildID()].ChannelID, "Queue ended.")
		}
		return true
	}
	return false
}
