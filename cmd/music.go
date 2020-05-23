package cmd

import (
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/music"
	"github.com/dustin/go-humanize"
	"time"
)

func PauseCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":information_source: There is currently no active connection.")
		return
	}

	conn.Lock()
	conn.StreamingSession.SetPaused(true)
	conn.Unlock()

	ctx.Reply(":pause_button: Paused music.")
}

func ResumeCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":information_source: There is currently no active connection.")
		return
	}

	conn.Lock()
	conn.StreamingSession.SetPaused(false)
	conn.Unlock()

	ctx.Reply(":arrow_forward: Resumed music.")
}

func SkipCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":information_source: There is currently no active connection.")
		return
	}

	conn.Lock()
	conn.Done <- errors.New("skip")
	conn.Unlock()
}

func StopCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":stop_button: There is currently no active connection.")
		return
	}

	conn.Lock()
	conn.Done <- errors.New("stop")
	conn.Unlock()
}

func NowPlayingCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The np command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":stop_button: There is currently no active connection.")
		return
	}

	conn.Lock()
	queueLen := conn.Queue.Len()
	if queueLen == 0 {
		ctx.Reply(":information_source: The current queue is empty.")
		conn.Unlock()
		return
	}

	firstTrack := conn.Queue.Front().Value.(*music.SongResponse)

	duration := "unknown"
	if !firstTrack.Duration.IsZero() {
		secs := int(firstTrack.Duration.Float64)
		duration = fmt.Sprintf("%d:%02d", secs/60, secs%60)
	}

	uploadDate := "Upload date unknown"
	if !firstTrack.UploadDate.IsZero() {
		t, err := time.Parse("20060102", firstTrack.UploadDate.String)
		if err == nil {
			uploadDate = "Uploaded " + humanize.Time(t)
		}
	}

	nowPlaying := &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("Duration %s - %s", duration, uploadDate),
				Value:  ":musical_note: " + firstTrack.Title,
				Inline: false,
			},
		},
		Title: "Now Playing",
	}

	if firstTrack.Thumbnail != "" {
		nowPlaying.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: firstTrack.Thumbnail,
		}
	}

	conn.Unlock()

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, nowPlaying)
}

func PlayCommand(ctx *exrouter.Context) {
	music.AddToQueue(ctx)
}

func QueueCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The queue command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":stop_button: There is currently no active connection.")
		return
	}

	conn.Lock()
	queueLen := conn.Queue.Len()
	if queueLen == 0 {
		ctx.Reply(":information_source: The current queue is empty.")
		conn.Unlock()
		return
	}

	queueList := &discordgo.MessageEmbed{
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{},
		Title:  "Music Queue",
	}

	firstTrack := conn.Queue.Front().Value.(*music.SongResponse)

	if firstTrack.Thumbnail != "" {
		queueList.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: firstTrack.Thumbnail,
		}
	}

	n := 0
	if queueLen < 10 {
		n = queueLen
	} else {
		n = 10
	}

	i := 0
	for el := conn.Queue.Front(); el != nil; el = el.Next() {
		if i >= n {
			break
		}

		currTrack := el.Value.(*music.SongResponse)

		title := currTrack.Title
		if i == 0 {
			title = ":arrow_forward: " + title
		} else {
			title = ":arrow_up: " + title
		}

		duration := "unknown"
		if !currTrack.Duration.IsZero() {
			secs := int(currTrack.Duration.Float64)
			duration = fmt.Sprintf("%d:%02d", secs/60, secs%60)
		}

		uploadDate := "Upload date unknown"
		if !currTrack.UploadDate.IsZero() {
			t, err := time.Parse("20060102", currTrack.UploadDate.String)
			if err == nil {
				uploadDate = "Uploaded " + humanize.Time(t)
			}
		}

		queueList.Fields = append(queueList.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Duration %s - %s", duration, uploadDate),
			Value:  title,
			Inline: false,
		})

		i++
	}

	conn.Unlock()

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, queueList)
}

func ClearCommand(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The clear command cannot be used in DMs.")
		return
	}

	conn, ok := music.ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ctx.Reply(":stop_button: There is currently no active connection.")
		return
	}

	conn.Lock()

	queueLen := conn.Queue.Len()
	if queueLen <= 1 {
		ctx.Reply(":information_source: No tracks to clear.")
		conn.Unlock()
		return
	}

	for conn.Queue.Len() != 1 {
		conn.Queue.Remove(conn.Queue.Back())
	}

	conn.Unlock()
	ctx.Reply(":put_litter_in_its_place: Queue cleared.")
}

func YoutubeCommand(ctx *exrouter.Context) {
	music.AddToQueueYT(ctx)
}
