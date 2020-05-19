package cmd

import (
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"github.com/foxbot/gavalink"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func musicPreprocessor(ctx *exrouter.Context) error {
	generalError := errors.New("setup error")

	// Check DMs
	if ctx.Msg.GuildID == "" {
		_, _ = ctx.Reply("Cannot play music in DMs.")
		return generalError
	}

	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)

	if err != nil {
		_, _ = ctx.Reply("An error has occurred when fetching information about your server.")
		log.Printf("music: %s", err)
		return generalError
	}

	if len(guild.VoiceStates) == 0 {
		_, _ = ctx.Reply("You are not in a voice channel.")
		return generalError
	}

	var state *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Msg.Author.ID {
			state = v
			break
		}
	}

	if state == nil {
		_, _ = ctx.Reply("You are not in a voice channel.")
		return generalError
	}

	var botState *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Ses.State.User.ID {
			botState = v
			break
		}
	}

	if botState == nil || music.AudioPlayers[guild.ID] == nil {
		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
		return generalError
	}

	if botState.ChannelID != state.ChannelID {
		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
		return generalError
	}

	return nil
}

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	// Return usage
	if query == "" {
		_, _ = ctx.Reply(fmt.Sprintf("Usage: `%splay <song url or playlist url>`", prefix))
		return
	}

	// Check DMs
	if ctx.Msg.GuildID == "" {
		_, _ = ctx.Reply("Cannot play music in DMs.")
		return
	}

	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)

	if err != nil {
		_, _ = ctx.Reply("An error has occurred when fetching information about your server.")
		log.Printf("music: %s", err)
		return
	}

	if len(guild.VoiceStates) == 0 {
		_, _ = ctx.Reply("You are not in a voice channel.")
		return
	}

	var state *discordgo.VoiceState

	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Msg.Author.ID {
			state = v
			break
		}
	}

	if state == nil {
		_, _ = ctx.Reply("You are not in a voice channel.")
		return
	}

	// Join voice
	if err := ctx.Ses.ChannelVoiceJoinManual(ctx.Msg.GuildID, state.ChannelID, false, false); err != nil {
		log.Printf("music: failed to join channel: %s", err)
		_, _ = ctx.Reply("Failed to join the voice channel.")
		return
	}

	// Wait for 5 seconds
	for i := 0; music.AudioPlayers[ctx.Msg.GuildID] == nil; i++ {
		if i > 4 {
			_, _ = ctx.Reply("Connection failed.")
			return
		}
		time.Sleep(time.Second * 1)
	}

	tracks, err := music.AudioNode.LoadTracks(url.QueryEscape(query))

	if err != nil || len(tracks.Tracks) == 0 {
		log.Printf("music: cannot load track(s): %s", query)
		_, _ = ctx.Reply("Couldn't find results for your query.")
		return
	}

	if len(tracks.Tracks) == 1 {
		_, _ = ctx.Reply(fmt.Sprintf("Added song: %s", tracks.Tracks[0].Info.Title))
	} else if len(tracks.Tracks) > 1 {
		_, _ = ctx.Reply(fmt.Sprintf("Added playlist: %s", tracks.PlaylistInfo.Name))
	}

	for _, val := range tracks.Tracks {
		if !music.AudioPlayers[ctx.Msg.GuildID].Playing {
			if err = music.AudioPlayers[ctx.Msg.GuildID].Player.Play(val.Data); err != nil {
				_, _ = ctx.Reply(fmt.Sprintf("Failed to play: %s", val.Info.Title))
			} else {
				music.AudioPlayers[ctx.Msg.GuildID].Playing = true
				music.AudioPlayers[ctx.Msg.GuildID].ChannelID = ctx.Msg.ChannelID
				music.AudioPlayers[ctx.Msg.GuildID].CurrTrack = val
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
					Color:       config.EmbedColor,
					Title:       "Now playing",
					Description: val.Info.Title,
					URL:         val.Info.URI,
				})
			}
		} else {
			music.AudioPlayers[ctx.Msg.GuildID].Queue = append(music.AudioPlayers[ctx.Msg.GuildID].Queue, val)
		}
	}
}

func LeaveCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	music.Disconnect(ctx)
}

func ResumeCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	if !music.AudioPlayers[ctx.Msg.GuildID].Player.Paused() {
		_, _ = ctx.Reply("Music is currently not paused.")
		return
	}

	if err := music.AudioPlayers[ctx.Msg.GuildID].Player.Pause(false); err != nil {
		log.Printf("music: resume fail: %s", err)
		_, _ = ctx.Reply("Failed to resume.")
		return
	}

	_, _ = ctx.Reply("Resumed music.")
}

func PauseCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	if music.AudioPlayers[ctx.Msg.GuildID].Player.Paused() {
		_, _ = ctx.Reply("Music is currently paused.")
		return
	}

	if err := music.AudioPlayers[ctx.Msg.GuildID].Player.Pause(true); err != nil {
		log.Printf("music: pause fail: %s", err)
		_, _ = ctx.Reply("Failed to pause.")
		return
	}

	_, _ = ctx.Reply("Paused music.")
}

func StopCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	if music.AudioPlayers[ctx.Msg.GuildID].Player.Position() == 0 {
		_, _ = ctx.Reply("Music is currently not playing.")
		return
	}

	if err := music.AudioPlayers[ctx.Msg.GuildID].Player.Stop(); err != nil {
		log.Printf("music: stop fail: %s", err)
		_, _ = ctx.Reply("Failed to stop.")
		return
	}

	music.AudioPlayers[ctx.Msg.GuildID].Queue = make([]gavalink.Track, 0)

	_, _ = ctx.Reply("Stopped music and cleared queue.")
}

func SkipCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	if err := music.AudioPlayers[ctx.Msg.GuildID].Player.Stop(); err != nil {
		log.Printf("music: stop fail: %s", err)
		_, _ = ctx.Reply("Failed to stop current track.")
		return
	}

	_, _ = ctx.Reply("Skipped current track.")
}

func ClearCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	music.AudioPlayers[ctx.Msg.GuildID].Queue = make([]gavalink.Track, 0)

	_, _ = ctx.Reply("Cleared current queue.")
}

func ShuffleCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	rand.Shuffle(len(music.AudioPlayers[ctx.Msg.GuildID].Queue),
		func(i, j int) { music.AudioPlayers[ctx.Msg.GuildID].Queue[i], music.AudioPlayers[ctx.Msg.GuildID].Queue[j] = music.AudioPlayers[ctx.Msg.GuildID].Queue[j], music.AudioPlayers[ctx.Msg.GuildID].Queue[i] })

	_, _ = ctx.Reply("Shuffled current queue.")
}

func QueueCommand(ctx *exrouter.Context) {
	if musicPreprocessor(ctx) != nil {
		return
	}

	_, _ = ctx.Reply(music.AudioPlayers[ctx.Msg.GuildID].Queue)
}
