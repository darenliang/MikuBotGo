// Code taken from ducc/GoMusicBot
package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"github.com/dustin/go-humanize"
	"log"
	"strings"
)

func StopCommand(ctx *exrouter.Context) {
	channel, err := ctx.Ses.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		channel, err = ctx.Ses.Channel(ctx.Msg.ChannelID)
		if err != nil {
			ctx.Reply(":cry: An error has occurred.")
			log.Print("music: channel get fail")
			return
		}
		ctx.Ses.State.ChannelAdd(channel)
	}

	conn, ok := music.MusicConnections[channel.GuildID]
	if !ok {
		ctx.Reply("No music is currently playing.")
		return
	}

	conn.Close()

	ctx.Reply(":octagonal_sign: Stopped music playback.")
}

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf("Usage: `%splay <url>`", prefix))
		return
	}

	channel, err := ctx.Ses.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		channel, err = ctx.Ses.Channel(ctx.Msg.ChannelID)
		if err != nil {
			log.Print("music: channel get fail")
			ctx.Reply(":cry: An error has occurred.")
			return
		}
		ctx.Ses.State.ChannelAdd(channel)
	}

	guild, err := ctx.Ses.State.Guild(channel.GuildID)
	if err != nil {
		guild, err = ctx.Ses.Guild(channel.GuildID)
		if err != nil {
			log.Print("music: guild get fail")
			ctx.Reply(":cry: An error has occurred.")
			return
		}
		ctx.Ses.State.GuildAdd(guild)
	}

	musicChannelID := ""
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == ctx.Msg.Author.ID && musicChannelID == "" {
			musicChannelID = voiceState.ChannelID
		}
	}

	conn, ok := music.MusicConnections[channel.GuildID]
	if ok {
		if musicChannelID != conn.ChannelID {
			ctx.Reply(":information_source: Please join a voice channel.")
			return
		}
		go music.AddToQueue(ctx, conn, query)
	} else {
		if musicChannelID == "" {
			ctx.Reply(":information_source: Please join a voice channel.")
			return
		}

		go music.PlaySong(ctx, channel.GuildID, musicChannelID, query)
	}
}

func QueueCommand(ctx *exrouter.Context) {
	channel, err := ctx.Ses.State.Channel(ctx.Msg.ChannelID)
	if err != nil {
		channel, err = ctx.Ses.Channel(ctx.Msg.ChannelID)
		if err != nil {
			log.Print("music: channel get fail")
			ctx.Reply(":cry: An error has occurred.")
			return
		}
		ctx.Ses.State.ChannelAdd(channel)
	}

	conn, ok := music.MusicConnections[channel.GuildID]
	if !ok {
		ctx.Reply(":mute: No music is currently playing.")
		return
	}

	queueLen := len(conn.Queue)
	if queueLen == 0 {
		ctx.Reply(":information_source: The current queue is empty.")
		return
	}

	queueList := &discordgo.MessageEmbed{
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{},
		Title:  "Music Queue",
	}

	n := 0
	if queueLen < 10 {
		n = queueLen
	} else {
		n = 10
	}

	for i := 0; i < n; i++ {
		queueItem := conn.Queue[i]
		title := queueItem.Info.Title
		if i == 0 {
			title = ":arrow_forward: " + title
		}
		queueList.Fields = append(queueList.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Duration %d:%02d:%02d - Uploaded %s", int(queueItem.Info.Duration.Hours()), int(queueItem.Info.Duration.Minutes())%60, int(queueItem.Info.Duration.Seconds())%60, humanize.Time(queueItem.Info.DatePublished)),
			Value:  title,
			Inline: false,
		})
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, queueList)
}
