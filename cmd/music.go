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
	"sync"
	"time"
)

func PauseCommand(ctx *exrouter.Context) {
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

	fin, err := conn.Stream.Finished()
	if fin {
		ctx.Reply(":mute: No music is currently playing.")
		return
	}

	if conn.Stream.Paused() {
		ctx.Reply(":information_source: Music is currently paused.")
		return
	}

	conn.Stream.SetPaused(true)
	ctx.Reply(":pause_button: Paused music.")
}

func ResumeCommand(ctx *exrouter.Context) {
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

	fin, err := conn.Stream.Finished()
	if fin {
		ctx.Reply(":mute: No music is currently playing.")
		return
	}

	if !conn.Stream.Paused() {
		ctx.Reply(":information_source: Music is currently playing.")
		return
	}

	conn.Stream.SetPaused(false)
	ctx.Reply(":arrow_forward: Resumed music.")
}

func SkipCommand(ctx *exrouter.Context) {
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

	fin, err := conn.Stream.Finished()
	if fin {
		ctx.Reply(":mute: No music is currently playing.")
		return
	}

	conn.Done <- nil
	ctx.Reply(":track_next: Skipped current track.")
}

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
		ctx.Reply(":mute: No music is currently playing.")
		return
	}

	conn.Close()

	ctx.Reply(":octagonal_sign: Stopped music playback. Queue cleared.")
}

func NowPlayingCommand(ctx *exrouter.Context) {
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
		ctx.Reply(":mute: No music is playing.")
		return
	}

	queueLen := len(conn.Queue)
	if queueLen == 0 {
		ctx.Reply(":information_source: The current queue is empty.")
		return
	}

	duration := "unknown"
	if !conn.Queue[0].Duration.IsZero() {
		secs := int(conn.Queue[0].Duration.Float64)
		duration = fmt.Sprintf("%d:%02d", secs/60, secs%60)
	}

	uploadDate := "Upload date unknown"
	if !conn.Queue[0].UploadDate.IsZero() {
		t, err := time.Parse("20060102", conn.Queue[0].UploadDate.String)
		if err == nil {
			uploadDate = "Uploaded " + humanize.Time(t)
		}
	}

	// TODO: Fix possible dangerous mutation to queue
	nowPlaying := &discordgo.MessageEmbed{
		Color: config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   fmt.Sprintf("Duration %s - %s", duration, uploadDate),
				Value:  ":musical_note: " + conn.Queue[0].Title,
				Inline: false,
			},
		},
		Title: "Now Playing",
	}

	if conn.Queue[0].Thumbnail != "" {
		nowPlaying.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: conn.Queue[0].Thumbnail,
		}
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, nowPlaying)
}

func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%splay <url or playlist url>`", prefix))
		return
	}

	if music.MusicCooldown[ctx.Msg.GuildID].IsZero() {
		music.MusicCooldown[ctx.Msg.GuildID] = time.Now()
	} else {
		elapsed := time.Since(music.MusicCooldown[ctx.Msg.GuildID])
		if elapsed.Seconds() <= music.MUSICCOOLDOWN {
			ctx.Reply(fmt.Sprintf(":timer: You need to wait %d seconds before your cooldown expires.", music.MUSICCOOLDOWN-int(elapsed.Seconds())))
			return
		}
		music.MusicCooldown[ctx.Msg.GuildID] = time.Now()
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

	botInVoice := false
	musicChannelID := ""
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == ctx.Msg.Author.ID && musicChannelID == "" {
			musicChannelID = voiceState.ChannelID
		} else if voiceState.UserID == ctx.Ses.State.User.ID {
			botInVoice = true
		}
	}

	if !botInVoice {
		delete(music.MusicConnections, guild.ID)
	}

	playReadyStatus, _, err := music.GetPlayReadyData(musicChannelID, channel.GuildID)
	if err != nil {
		ctx.Reply(":information_source: Please join a joice channel.")
		return
	}

	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	t, inp, err := music.Youtube{}.YoutubeDLLink(query)

	if err != nil {
		log.Printf("music: ytdl fail: %s", query)
		ctx.Reply(":cry: Cannot process your query.")
		return
	}

	switch t {
	case music.ERRORTYPE:
		{
			log.Printf("music: add song error: %s", query)
			ctx.Reply(":cry: An error occurred when getting song info.")
			return
		}
	case music.SONGTYPE:
		{
			song, err := music.Youtube{}.GetSong(*inp)

			if err != nil {
				ctx.Reply(":cry: Cannot find music.")
				log.Printf("music: add song video fail: %s", query)
				return
			}
			if playReadyStatus {
				go music.PlaySong(ctx, channel.GuildID, musicChannelID, song)
			} else {
				_, conn, err := music.GetPlayReadyData(musicChannelID, channel.GuildID)
				if err != nil {
					ctx.Reply(":warning: Voice connection error.")
					log.Printf("music: voice connection state fail")
					return
				}
				go music.AddToQueue(ctx, conn, song)
			}
		}
	case music.PLAYLISTTYPE:
		{
			songs, err := music.Youtube{}.GetPlaylist(*inp)
			if err != nil {
				ctx.Reply(":cry: An error occurred when getting playlist info.")
				log.Printf("music: add song playlist fail: %s", query)
				return
			}
			if len(*songs) > 10 {
				ctx.Reply(":information_source: Playlist is greater than 10 tracks. Only the first 10 tracks will be used.")
			}

			for idx, v := range *songs {
				if idx == 10 {
					break
				}
				url := v.URL
				_, i, err := music.Youtube{}.YoutubeDLLink(url)
				if err != nil {
					log.Printf("music: add song playlist song fail: %s", url)
					continue
				}
				song, err := music.Youtube{}.GetSong(*i)
				if err != nil {
					log.Printf("music: add song playlist song fail: %s", url)
					continue
				}
				if playReadyStatus {
					playReadyStatus = false
					go music.PlaySong(ctx, channel.GuildID, musicChannelID, song)
				} else {
					_, conn, err := music.GetPlayReadyData(musicChannelID, channel.GuildID)
					if err != nil {
						ctx.Reply(":warning: Voice connection error.")
						log.Printf("music: voice connection state fail")
						return
					}
					go music.AddToQueue(ctx, conn, song)
				}
			}
		}
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
		ctx.Reply(":mute: No music is playing.")
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

	if conn.Queue[0].Thumbnail != "" {
		queueList.Thumbnail = &discordgo.MessageEmbedThumbnail{
			URL: conn.Queue[0].Thumbnail,
		}
	}

	n := 0
	if queueLen < 10 {
		n = queueLen
	} else {
		n = 10
	}

	// TODO: Fix possible dangerous mutation to queue
	for i := 0; i < n; i++ {
		queueItem := conn.Queue[i]
		title := queueItem.Title
		if i == 0 {
			title = ":arrow_forward: " + title
		} else {
			title = ":arrow_up: " + title
		}

		duration := "unknown"
		if !queueItem.Duration.IsZero() {
			secs := int(queueItem.Duration.Float64)
			duration = fmt.Sprintf("%d:%02d", secs/60, secs%60)
		}

		uploadDate := "Upload date unknown"
		if !queueItem.UploadDate.IsZero() {
			t, err := time.Parse("20060102", queueItem.UploadDate.String)
			if err == nil {
				uploadDate = "Uploaded " + humanize.Time(t)
			}
		}

		queueList.Fields = append(queueList.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("Duration %s - %s", duration, uploadDate),
			Value:  title,
			Inline: false,
		})
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, queueList)
}

func ClearCommand(ctx *exrouter.Context) {
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
		ctx.Reply(":mute: No music is playing.")
		return
	}

	queueLen := len(conn.Queue)
	if queueLen == 0 {
		ctx.Reply(":information_source: The current queue is empty.")
		return
	}
	if queueLen == 1 {
		ctx.Reply(":information_source: There are no songs currently queued.")
		return
	}

	conn.Mutex.Lock()
	conn.Queue = conn.Queue[:1]
	conn.Mutex.Unlock()

	ctx.Reply(":put_litter_in_its_place: Queue cleared.")
}

func YoutubeCommand(ctx *exrouter.Context) {
	var (
		lock     sync.RWMutex
		embedMsg *discordgo.Message
		callback = make(chan struct{})
		idx      int
	)

	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%syt <query>`", prefix))
		return
	}

	if music.MusicCooldown[ctx.Msg.GuildID].IsZero() {
		music.MusicCooldown[ctx.Msg.GuildID] = time.Now()
	} else {
		elapsed := time.Since(music.MusicCooldown[ctx.Msg.GuildID])
		if elapsed.Seconds() <= music.MUSICCOOLDOWN {
			ctx.Reply(fmt.Sprintf(":timer: You need to wait %d seconds before your cooldown expires.", music.MUSICCOOLDOWN-int(elapsed.Seconds())))
			return
		}
		music.MusicCooldown[ctx.Msg.GuildID] = time.Now()
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

	botInVoice := false
	musicChannelID := ""
	for _, voiceState := range guild.VoiceStates {
		if voiceState.UserID == ctx.Msg.Author.ID && musicChannelID == "" {
			musicChannelID = voiceState.ChannelID
		} else if voiceState.UserID == ctx.Ses.State.User.ID {
			botInVoice = true
		}
	}

	if !botInVoice {
		delete(music.MusicConnections, guild.ID)
	}

	playReadyStatus, _, err := music.GetPlayReadyData(musicChannelID, channel.GuildID)
	if err != nil {
		ctx.Reply(":information_source: Please join a joice channel.")
		return
	}

	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	_, inp, err := music.Youtube{}.YoutubeDLQuery(query)

	if err != nil {
		log.Printf("music: ytdl search fail: %s", query)
		ctx.Reply(":cry: An error occurred when searching.")
		return
	}

	songs, err := music.Youtube{}.GetPlaylist(*inp)
	if err != nil {
		ctx.Reply(":cry: An error occurred when searching.")
		log.Printf("music: ytdl search songs fail: %s", query)
		return
	}

	// TODO: Fix case where results are less than 4
	if len(*songs) < 4 {
		ctx.Reply(":cry: I'm not getting enough search results.")
		log.Printf("music: not enough search results: %s", query)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  0xc4302b,
		Description: fmt.Sprintf(
			":one: %s\n"+
				":two: %s\n"+
				":three: %s\n"+
				":four: %s",
			(*songs)[0].Title,
			(*songs)[1].Title,
			(*songs)[2].Title,
			(*songs)[3].Title,
		),
		Title: fmt.Sprintf("Youtube Search Results for %s", query),
	}

	// Be a bad boy and reuse emojis from a weird location
	embedMsg, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	for i := 0; i < 4; i++ {
		ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, embedMsg.ID, emojis[i])
	}

	// Delete after command is done
	defer ctx.Ses.ChannelMessageDelete(ctx.Msg.ChannelID, embedMsg.ID)

	defer ctx.Ses.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		lock.RLock()
		defer lock.RUnlock()

		idx = framework.Index(emojis, reaction.Emoji.Name)

		if reaction.MessageID == embedMsg.ID && reaction.UserID != ctx.Ses.State.User.ID && idx != -1 {
			callback <- struct{}{}
		}
	})()

	select {
	case <-callback:
		if !(0 <= idx && idx <= 3) {
			log.Printf("music: index fail: %d", idx)
			ctx.Reply(":cry: An unexpected error has occurred")
			return
		}

		_, i, err := music.Youtube{}.YoutubeDLLink((*songs)[idx].URL)

		if err != nil {
			log.Printf("music: add song fail: %s", (*songs)[idx].Title)
			ctx.Reply(":cry: Failed to play song.")
			return
		}

		song, err := music.Youtube{}.GetSong(*i)

		if err != nil {
			log.Printf("music: add song fail: %s", (*songs)[idx].Title)
			ctx.Reply(":cry: Failed to play song.")
			return
		}
		if playReadyStatus {
			go music.PlaySong(ctx, channel.GuildID, musicChannelID, song)
		} else {
			_, conn, err := music.GetPlayReadyData(musicChannelID, channel.GuildID)
			if err != nil {
				ctx.Reply(":warning: Voice connection error.")
				log.Printf("music: voice connection state fail")
				return
			}
			go music.AddToQueue(ctx, conn, song)
		}
	case <-time.After(config.Timeout * time.Second):
		ctx.Reply(":timer: Youtube search timed out.")
	}
}
