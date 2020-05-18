package cmd

import (
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/foxbot/gavalink"
	"log"
	"os"
	"strings"
	"time"
)

type GuildPlayer struct {
	Player *gavalink.Player
	Queue  []gavalink.Track
}

var (
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

func JoinChannel(ctx *exrouter.Context) (bool, error) {
	if ctx.Msg.GuildID == "" {
		_, _ = ctx.Reply("Cannot play music in DMs.")
		return false, errors.New("music in dms")
	}

	g, err := ctx.Ses.State.Guild(ctx.Msg.GuildID)

	if err != nil {
		return false, err
	}

	var as *discordgo.VoiceState
	var bs *discordgo.VoiceState
	var moveTo string

	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.Msg.Author.ID {
			as = vs
		} else if vs.UserID == ctx.Ses.State.User.ID {
			bs = vs
		}
	}

	if as == nil {
		_, err := ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You must be in a voice channel.")
		if err != nil {
			return false, err
		}
		return false, nil
	} else if bs != nil && bs.ChannelID != as.ChannelID {
		// TODO: check permissions
		moveTo = as.ChannelID
	} else if bs == nil {
		moveTo = as.ChannelID
	}

	if len(moveTo) > 0 {
		err = ctx.Ses.ChannelVoiceJoinManual(ctx.Msg.GuildID, moveTo, false, false)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if query == "" {
		_, _ = ctx.Reply(fmt.Sprintf("Usage: `%splay <url>`", prefix))
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

	if err := ctx.Ses.ChannelVoiceJoinManual(ctx.Msg.GuildID, state.ChannelID, false, false); err != nil {
		log.Printf("music: failed to join channel: %s", err)
		_, _ = ctx.Reply("Failed to join the voice channel.")
		return
	}

	for i := 0; AudioPlayers[ctx.Msg.GuildID] == nil; i++ {
		if i > 4 {
			_, _ = ctx.Reply("Connection failed.")
			return
		}
		time.Sleep(time.Second * 1)
	}

	tracks, err := AudioNode.LoadTracks(query)

	if err != nil || len(tracks.Tracks) == 0 {
		log.Printf("music: cannot load track(s): %s", query)
		_, _ = ctx.Reply("Couldn't find anything for your query.")
		return
	}

	playing := false
	for idx, val := range tracks.Tracks {
		if !playing {
			if err = AudioPlayers[ctx.Msg.GuildID].Player.Play(tracks.Tracks[idx].Data); err != nil {
				_, _ = ctx.Reply(fmt.Sprintf("Failed to play: %s", tracks.Tracks[idx].Info.Title))
			} else {
				playing = true
			}
		} else {
			AudioPlayers[ctx.Msg.GuildID].Queue = append(AudioPlayers[ctx.Msg.GuildID].Queue, val)
			_, _ = ctx.Reply(fmt.Sprintf("Added to Queue: %s", tracks.Tracks[idx].Info.Title))
		}
	}

	if playing {
		_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, &discordgo.MessageEmbed{
			Color:       config.EmbedColor,
			Title:       "Now playing",
			Description: tracks.Tracks[0].Info.Title,
			URL:         tracks.Tracks[0].Info.URI,
		})
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
