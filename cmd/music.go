package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"log"
	"net/url"
	"strings"
	"time"
)

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	// Usage
	if query == "" {
		_, _ = ctx.Reply(fmt.Sprintf("Usage: `%splay <song url or playlist url>`", prefix))
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
		_, _ = ctx.Reply("Couldn't find anything for your query.")
		return
	}

	if len(tracks.Tracks) == 1 {
		_, _ = ctx.Reply(fmt.Sprintf("Adding song: %s", tracks.Tracks[0].Info.Title))
	} else if len(tracks.Tracks) > 1 {
		_, _ = ctx.Reply(fmt.Sprintf("Added playlist: %s", tracks.PlaylistInfo.Name))
	}

	for idx, val := range tracks.Tracks {
		if !music.AudioPlayers[ctx.Msg.GuildID].Playing {
			if err = music.AudioPlayers[ctx.Msg.GuildID].Player.Play(tracks.Tracks[idx].Data); err != nil {
				_, _ = ctx.Reply(fmt.Sprintf("Failed to play: %s", tracks.Tracks[idx].Info.Title))
			} else {
				music.AudioPlayers[ctx.Msg.GuildID].Playing = true
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
					Color:       config.EmbedColor,
					Title:       "Now playing",
					Description: tracks.Tracks[idx].Info.Title,
					URL:         tracks.Tracks[idx].Info.URI,
				})
			}
		} else {
			music.AudioPlayers[ctx.Msg.GuildID].Queue = append(music.AudioPlayers[ctx.Msg.GuildID].Queue, val)
		}
	}
}

//
// func ResumeCommand(ctx *exrouter.Context) {
// 	if ctx.Msg.GuildID == "" {
// 		_, _ = ctx.Reply("Cannot play music in DMs.")
// 		return
// 	}
//
// 	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)
//
// 	if err != nil {
// 		_, _ = ctx.Reply("An error has occurred when fetching information about your server.")
// 		log.Printf("music: %s", err)
// 		return
// 	}
//
// 	if len(guild.VoiceStates) == 0 {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var state *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Msg.Author.ID {
// 			state = v
// 			break
// 		}
// 	}
//
// 	if state == nil {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var botState *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Ses.State.User.ID {
// 			botState = v
// 			break
// 		}
// 	}
//
// 	if botState == nil || AudioPlayers[guild.ID] == nil {
// 		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
// 		return
// 	}
//
// 	if botState.ChannelID != state.ChannelID {
// 		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
// 		return
// 	}
//
// 	if !AudioPlayers[guild.ID].Player.Paused() {
// 		_, _ = ctx.Reply("Music is currently not paused.")
// 		return
// 	}
//
// 	if err := AudioPlayers[guild.ID].Player.Pause(false); err != nil {
// 		log.Printf("music: resume fail: %s", err)
// 		_, _ = ctx.Reply("Failed to resume.")
// 		return
// 	}
//
// 	_, _ = ctx.Reply("Resumed music.")
// }
//
// func PauseCommand(ctx *exrouter.Context) {
// 	if ctx.Msg.GuildID == "" {
// 		_, _ = ctx.Reply("Cannot play music in DMs.")
// 		return
// 	}
//
// 	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)
//
// 	if err != nil {
// 		_, _ = ctx.Reply("An error has occurred when fetching information about your server.")
// 		log.Printf("music: %s", err)
// 		return
// 	}
//
// 	if len(guild.VoiceStates) == 0 {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var state *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Msg.Author.ID {
// 			state = v
// 			break
// 		}
// 	}
//
// 	if state == nil {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var botState *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Ses.State.User.ID {
// 			botState = v
// 			break
// 		}
// 	}
//
// 	if botState == nil || AudioPlayers[guild.ID] == nil {
// 		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
// 		return
// 	}
//
// 	if botState.ChannelID != state.ChannelID {
// 		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
// 		return
// 	}
//
// 	if AudioPlayers[guild.ID].Player.Paused() {
// 		_, _ = ctx.Reply("Music is currently paused.")
// 		return
// 	}
//
// 	if err := AudioPlayers[guild.ID].Player.Pause(true); err != nil {
// 		log.Printf("music: pause fail: %s", err)
// 		_, _ = ctx.Reply("Failed to pause.")
// 		return
// 	}
//
// 	_, _ = ctx.Reply("Paused music.")
// }
//
// func StopCommand(ctx *exrouter.Context) {
// 	if ctx.Msg.GuildID == "" {
// 		_, _ = ctx.Reply("Cannot play music in DMs.")
// 		return
// 	}
//
// 	guild, err := ctx.Ses.Guild(ctx.Msg.GuildID)
//
// 	if err != nil {
// 		_, _ = ctx.Reply("An error has occurred when fetching information about your server.")
// 		log.Printf("music: %s", err)
// 		return
// 	}
//
// 	if len(guild.VoiceStates) == 0 {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var state *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Msg.Author.ID {
// 			state = v
// 			break
// 		}
// 	}
//
// 	if state == nil {
// 		_, _ = ctx.Reply("You are not in a voice channel.")
// 		return
// 	}
//
// 	var botState *discordgo.VoiceState
// 	for _, v := range guild.VoiceStates {
// 		if v.UserID == ctx.Ses.State.User.ID {
// 			botState = v
// 			break
// 		}
// 	}
//
// 	if botState == nil || AudioPlayers[guild.ID] == nil {
// 		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
// 		return
// 	}
//
// 	if botState.ChannelID != state.ChannelID {
// 		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
// 		return
// 	}
//
// 	if AudioPlayers[guild.ID].Player.Position() == 0 {
// 		_, _ = ctx.Reply("Music is currently not playing.")
// 		return
// 	}
//
// 	AudioPlayers[guild.ID].Player
//
// 	if err := AudioPlayers[guild.ID].Player.Stop(); err != nil {
// 		log.Printf("music: stop fail: %s", err)
// 		_, _ = ctx.Reply("Failed to stop.")
// 		return
// 	}
//
// 	ctx.Ses.VoiceConnections[guild.ID].Close()
//
// 	_, _ = ctx.Reply("Stopped music.")
// }
