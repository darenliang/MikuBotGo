package music

import (
	"fmt"
	"github.com/foxbot/gavalink"
)

const stuckTimeout = 10 * 1000

type EventHandler struct{}

func (m EventHandler) OnTrackEnd(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player) {
		Disconnect(player.GuildID())
		return nil
	}

	PlayNextTrack(player)
	return nil
}

func (m EventHandler) OnTrackException(player *gavalink.Player, track string, reason string) error {
	if IsEmptyQueue(player) {
		Disconnect(player.GuildID())
		return nil
	}

	PlayNextTrack(player)
	return nil
}

func (m EventHandler) OnTrackStuck(player *gavalink.Player, track string, threshold int) error {
	if IsEmptyQueue(player) {
		return nil
	}

	if threshold > stuckTimeout {
		PlayNextTrack(player)
	}
	return nil
}

func PlayNextTrack(player *gavalink.Player) {
	guildPlayer := AudioPlayers[player.GuildID()]
	queue := guildPlayer.Queue
	for queue.Len() != 0 {
		track := queue.Remove(queue.Front()).(gavalink.Track)
		if err := player.Play(track.Data); err != nil {
			Session.ChannelMessageSend(guildPlayer.ChannelID, fmt.Sprintf("Failed to play: %s", track.Info.Title))
			continue
		}
		break
	}
}

func IsEmptyQueue(player *gavalink.Player) bool {
	return AudioPlayers[player.GuildID()].Queue.Len() == 0
}
