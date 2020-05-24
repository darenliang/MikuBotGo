package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"github.com/foxbot/gavalink"
	"log"
	"net/url"
	"strings"
	"time"
)

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if query == "" {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%splay <song url or playlist url>`", prefix))
		return
	}

	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The gif command cannot be used in DMs.")
		return
	}

	guild, err := ctx.Ses.State.Guild(ctx.Msg.GuildID)
	if err != nil {
		guild, err = ctx.Ses.Guild(ctx.Msg.GuildID)
		if err != nil {
			ctx.Reply(":cry: An extremely rare error has occurred.")
			return
		}
		ctx.Ses.State.GuildAdd(guild)
	}

	var state *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Msg.Author.ID {
			state = v
			break
		}
	}

	if state == nil {
		ctx.Reply("You are not in a voice channel.")
		return
	}

	if err := ctx.Ses.ChannelVoiceJoinManual(ctx.Msg.GuildID, state.ChannelID, false, false); err != nil {
		log.Printf("music: failed to join channel: %s", err)
		_, _ = ctx.Reply(":cry: Failed to join the voice channel.")
		return
	}

	var player *gavalink.Player

	// Wait for 10 seconds
	for i := 0; ; i++ {
		if i > 4 {
			ctx.Reply(":cry: Connection failed.")
			return
		}
		if player, err = music.AudioLavalink.GetPlayer(ctx.Msg.GuildID); err == nil {
			break
		}
		time.Sleep(time.Second * 1)
	}

	node, err := music.AudioLavalink.BestNode()
	if err != nil {
		ctx.Reply(":cry: Failed to find optimal music node.")
		return
	}

	tracks, err := node.LoadTracks(url.QueryEscape(query))
	if err != nil || len(tracks.Tracks) == 0 {
		ctx.Reply(":cry: Cannot process your query.")
		return
	}

	switch tracks.Type {
	case gavalink.TrackLoaded:
		{
			if player.Position() == 0 {
				if err = player.Play(tracks.Tracks[0].Data); err != nil {
					_, _ = ctx.Reply(fmt.Sprintf(":cry: Failed to play: %s", tracks.Tracks[0].Info.Title))
					return
				}
				music.AudioPlayers[ctx.Msg.GuildID].ChannelID = ctx.Msg.ChannelID
				ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
					Color:       config.EmbedColor,
					Title:       "Now Playing",
					Description: tracks.Tracks[0].Info.Title,
					URL:         tracks.Tracks[0].Info.URI,
				})
			} else {
				ctx.Reply(fmt.Sprintf(":white_check_mark: Added track: %s", tracks.Tracks[0].Info.Title))
			}
			music.AudioPlayers[ctx.Msg.GuildID].Queue.PushBack(tracks.Tracks[0])
			fmt.Println(music.AudioPlayers[ctx.Msg.GuildID].Queue.Len())
		}
	case gavalink.PlaylistLoaded:
		{
			ctx.Reply(fmt.Sprintf(":white_check_mark: Added playlist: %s", tracks.PlaylistInfo.Name))
			stopped := player.Position() == 0
			for _, track := range tracks.Tracks {
				if stopped {
					if err = player.Play(track.Data); err != nil {
						_, _ = ctx.Reply(fmt.Sprintf(":cry: Failed to play: %s", track.Info.Title))
						continue
					}
					stopped = false
					music.AudioPlayers[ctx.Msg.GuildID].ChannelID = ctx.Msg.ChannelID
					ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
						Color:       config.EmbedColor,
						Title:       "Now Playing",
						Description: track.Info.Title,
						URL:         track.Info.URI,
					})
				}
			}
			music.AudioPlayers[ctx.Msg.GuildID].Queue.PushBack(tracks.Tracks[0])
		}
	case gavalink.LoadFailed:
		{
			ctx.Reply(":cry: Cannot process your query.")
			return
		}
	case gavalink.NoMatches:
		{
			ctx.Reply(":cry: No results found.")
			return
		}
	default:
		{
			// TODO fix this
			ctx.Reply(":information_source: You've a weird query over there. Could you provide an url instead?")
			return
		}
	}
}

func PauseCommand(ctx *exrouter.Context) {}

func SkipCommand(ctx *exrouter.Context) {}

func StopCommand(ctx *exrouter.Context) {}

func NowPlayingCommand(ctx *exrouter.Context) {}

func QueueCommand(ctx *exrouter.Context) {}

func ClearCommand(ctx *exrouter.Context) {}
