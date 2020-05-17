package cmd

import (
	"bytes"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/foxbot/gavalink"
	"github.com/robfig/cron"
	"log"
	"os"
	"strings"
	"time"
)

var (
	AudioLavalink *gavalink.Lavalink
	AudioNode     *gavalink.Node
	AudioPlayers  map[string]*gavalink.Player
	lavalinkRest  string
	lavalinkWS    string
	lavalinkPass  string
)

func init() {
	lavalinkRest = os.Getenv("LAVALINK_REST")
	lavalinkWS = os.Getenv("LAVALINK_WS")
	lavalinkPass = os.Getenv("LAVALINK_PASS")
}

func AudioInit() {
	var (
		buf      = bytes.NewBuffer([]byte{})
		c        = cron.New()
		prevcont = buf.String()
	)

	gavalink.Log = log.New(buf, "[gavalink] ", 0)

	_ = c.AddFunc("@every 500ms", func() {
		if buf.String() == prevcont {
			return
		}
		buf.Reset()
	})

	c.Start()

	AudioLavalink = gavalink.NewLavalink("1", config.BotID)
	AudioPlayers = make(map[string]*gavalink.Player)

	err := AudioLavalink.AddNodes(gavalink.NodeConfig{
		REST:      "http://lavalink-discord.herokuapp.com:80",
		WebSocket: "ws://lavalink-discord.herokuapp.com:80",
		Password:  "youshallnotpass",
	})

	if err != nil {
		panic(err)
	}
}

func PlayCommand(ctx *exrouter.Context) {
	query := strings.TrimSpace(ctx.Args.After(1))

	if query == "" {
		_, _ = ctx.Reply("Please provide a query.")
		return
	}

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

	if err = AudioPlayers[ctx.Msg.GuildID].Play(tracks.Tracks[0].Data); err != nil {
		_, _ = ctx.Reply("Encountered playback error.")
		return
	}

	if tracks.Type != gavalink.TrackLoaded {
		_, _ = ctx.Reply("Cannot play playlists yet.")
		return
	}

	_, _ = ctx.ReplyEmbed(&discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{{
			Name:   "Title",
			Value:  tracks.Tracks[0].Info.Title,
			Inline: false,
		}, {
			Name:   "Author",
			Value:  tracks.Tracks[0].Info.Author,
			Inline: false,
		}},
		Title: "Now playing",
		URL:   tracks.Tracks[0].Info.URI,
	})
}

func ResumeCommand(ctx *exrouter.Context) {
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

	var botState *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Ses.State.User.ID {
			botState = v
			break
		}
	}

	if botState == nil || AudioPlayers[guild.ID] == nil {
		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
		return
	}

	if botState.ChannelID != state.ChannelID {
		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
		return
	}

	if !AudioPlayers[guild.ID].Paused() {
		_, _ = ctx.Reply("Music is currently not paused.")
		return
	}

	if err := AudioPlayers[guild.ID].Pause(false); err != nil {
		log.Printf("music: resume fail: %s", err)
		_, _ = ctx.Reply("Failed to resume.")
		return
	}

	_, _ = ctx.Reply("Resumed music.")
}

func PauseCommand(ctx *exrouter.Context) {
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

	var botState *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Ses.State.User.ID {
			botState = v
			break
		}
	}

	if botState == nil || AudioPlayers[guild.ID] == nil {
		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
		return
	}

	if botState.ChannelID != state.ChannelID {
		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
		return
	}

	if AudioPlayers[guild.ID].Paused() {
		_, _ = ctx.Reply("Music is currently paused.")
		return
	}

	if err := AudioPlayers[guild.ID].Pause(true); err != nil {
		log.Printf("music: pause fail: %s", err)
		_, _ = ctx.Reply("Failed to pause.")
		return
	}

	_, _ = ctx.Reply("Paused music.")
}

func StopCommand(ctx *exrouter.Context) {
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

	var botState *discordgo.VoiceState
	for _, v := range guild.VoiceStates {
		if v.UserID == ctx.Ses.State.User.ID {
			botState = v
			break
		}
	}

	if botState == nil || AudioPlayers[guild.ID] == nil {
		_, _ = ctx.Reply("The bot is currently not in a voice channel.")
		return
	}

	if botState.ChannelID != state.ChannelID {
		_, _ = ctx.Reply("The bot is currently not in the same voice channel.")
		return
	}

	if AudioPlayers[guild.ID].Position() == 0 {
		_, _ = ctx.Reply("Music is currently not playing.")
		return
	}

	if err := AudioPlayers[guild.ID].Stop(); err != nil {
		log.Printf("music: stop fail: %s", err)
		_, _ = ctx.Reply("Failed to stop.")
		return
	}

	ctx.Ses.VoiceConnections[guild.ID].Close()

	_, _ = ctx.Reply("Stopped music.")
}
