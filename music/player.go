package music

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/jonas747/dca"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const COOLDOWNTIME = 10

type Connections struct {
	*sync.RWMutex
	ConnectionMap map[string]*Connection
}

// MusicConnections maps a Guild ID to an associated voice connection.
var (
	ServerConnections = Connections{ConnectionMap: make(map[string]*Connection)}
	UserJoinError     = errors.New("you are not in a voice channel")
	Client            = http.Client{Timeout: time.Second * 10}
)

func AddToQueue(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%splay <url or playlist url>`", prefix))
		return
	}

	conn, ok := ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ServerConnections.ConnectionMap[ctx.Msg.GuildID] = &Connection{
			RWMutex:    new(sync.RWMutex),
			Done:       make(chan error),
			Queue:      list.New(),
			LastInvoke: time.Now(),
		}
		conn = ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	} else {
		elapsed := time.Since(conn.LastInvoke)
		if elapsed.Seconds() <= COOLDOWNTIME {
			ctx.Reply(fmt.Sprintf(":timer: You need to wait %d seconds before your cooldown expires.", COOLDOWNTIME-int(elapsed.Seconds())))
			return
		}
		conn.LastInvoke = time.Now()
	}

	guild, err := ctx.Ses.State.Guild(ctx.Msg.GuildID)
	if err != nil {
		guild, err = ctx.Ses.Guild(ctx.Msg.GuildID)
		if err != nil {
			log.Print("music: guild get fail")
			ctx.Reply(":cry: An error has occurred.")
			return
		}
		ctx.Ses.State.GuildAdd(guild)
	}

	conn.Lock()
	defer conn.Unlock()
	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)
	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	vc, err := CreateVoiceConnection(ctx, guild, conn)
	if err != nil {
		return
	}

	t, inp, err := Youtube{}.YoutubeDLLink(query)

	if err != nil {
		log.Printf("music: ytdl fail: %s", query)
		ctx.Reply(":cry: Cannot process your query.")
		return
	}

	switch t {
	case ERRORTYPE:
		{
			log.Printf("music: add song error: %s", query)
			ctx.Reply(":cry: An error occurred when getting song info.")
			return
		}
	case SONGTYPE:
		{
			song, err := Youtube{}.GetSong(*inp)

			if err != nil {
				ctx.Reply(":cry: Cannot find music.")
				log.Printf("music: add song video fail: %s", query)
				return
			}

			conn.Queue.PushBack(song)
			ctx.Reply(fmt.Sprintf(":white_check_mark: Added %s to the queue.", song.Title))

			if !conn.Playing {
				conn.VoiceConnection = vc
				conn.Playing = true
				conn.VoiceChannelID = vc.ChannelID
				go Play(ctx, conn, vc)
			}
		}
	case PLAYLISTTYPE:
		{
			songs, err := Youtube{}.GetPlaylist(*inp)
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
				_, i, err := Youtube{}.YoutubeDLLink(url)
				if err != nil {
					log.Printf("music: add song playlist song fail: %s", url)
					continue
				}

				song, err := Youtube{}.GetSong(*i)
				if err != nil {
					log.Printf("music: add song playlist song fail: %s", url)
					continue
				}

				conn.Queue.PushBack(song)
				ctx.Reply(fmt.Sprintf(":white_check_mark: Added %s to the queue.", song.Title))

				if !conn.Playing {
					conn.VoiceConnection = vc
					conn.Playing = true
					conn.VoiceChannelID = vc.ChannelID
					go Play(ctx, conn, vc)
				}
			}
		}
	}
}

func AddToQueueYT(ctx *exrouter.Context) {
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The resume command cannot be used in DMs.")
		return
	}

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

	conn, ok := ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	if !ok {
		ServerConnections.ConnectionMap[ctx.Msg.GuildID] = &Connection{
			RWMutex:    new(sync.RWMutex),
			Done:       make(chan error),
			Queue:      list.New(),
			LastInvoke: time.Now(),
		}
		conn = ServerConnections.ConnectionMap[ctx.Msg.GuildID]
	} else {
		elapsed := time.Since(conn.LastInvoke)
		if elapsed.Seconds() <= COOLDOWNTIME {
			ctx.Reply(fmt.Sprintf(":timer: You need to wait %d seconds before your cooldown expires.", COOLDOWNTIME-int(elapsed.Seconds())))
			return
		}
		conn.LastInvoke = time.Now()
	}

	guild, err := ctx.Ses.State.Guild(ctx.Msg.GuildID)
	if err != nil {
		guild, err = ctx.Ses.Guild(ctx.Msg.GuildID)
		if err != nil {
			log.Print("music: guild get fail")
			ctx.Reply(":cry: An error has occurred.")
			return
		}
		ctx.Ses.State.GuildAdd(guild)
	}

	conn.Lock()
	defer conn.Unlock()
	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)
	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	vc, err := CreateVoiceConnection(ctx, guild, conn)
	if err != nil {
		return
	}

	_, inp, err := Youtube{}.YoutubeDLQuery(query)

	if err != nil {
		log.Printf("music: ytdl search fail: %s", query)
		ctx.Reply(":cry: An error occurred when searching.")
		return
	}

	songs, err := Youtube{}.GetPlaylist(*inp)
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

	embedMsg, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	for i := 0; i < 4; i++ {
		ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, embedMsg.ID, config.SelectionEmojis[i])
	}

	// Delete after command is done
	defer ctx.Ses.ChannelMessageDelete(ctx.Msg.ChannelID, embedMsg.ID)

	defer ctx.Ses.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		lock.RLock()
		defer lock.RUnlock()

		idx = framework.Index(config.SelectionEmojis, reaction.Emoji.Name)

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

		_, i, err := Youtube{}.YoutubeDLLink((*songs)[idx].URL)

		if err != nil {
			log.Printf("music: add song fail: %s", (*songs)[idx].Title)
			ctx.Reply(":cry: Failed to play song.")
			return
		}

		song, err := Youtube{}.GetSong(*i)

		if err != nil {
			log.Printf("music: add song fail: %s", (*songs)[idx].Title)
			ctx.Reply(":cry: Failed to play song.")
			return
		}

		conn.Queue.PushBack(song)
		ctx.Reply(fmt.Sprintf(":white_check_mark: Added %s to the queue.", song.Title))

		if !conn.Playing {
			conn.VoiceConnection = vc
			conn.Playing = true
			conn.VoiceChannelID = vc.ChannelID
			go Play(ctx, conn, vc)
		}
	case <-time.After(config.Timeout * time.Second):
		ctx.Reply(":timer: Youtube search timed out.")
	}
}

func Play(ctx *exrouter.Context, conn *Connection, vc *discordgo.VoiceConnection) {
	if conn.Queue.Len() == 0 {
		conn.YoutubeCleanup()
		ctx.Reply(":information_source: Queue ended.")
		return
	}

	conn.Lock()
	song := conn.Queue.Front().Value.(*SongResponse)
	encSesh, err := dca.EncodeFile(song.URL, dca.StdEncodeOptions)
	if err != nil {
		log.Printf("player: song failed to play %s", song.Title)
		ctx.Reply(":information_source: Song failed to play.")
		conn.YoutubeCleanup()
		conn.Unlock()
		return
	}
	defer encSesh.Cleanup()

	conn.StreamingSession = dca.NewStream(encSesh, vc, conn.Done)
	conn.Unlock()

	ctx.Reply(fmt.Sprintf(":arrow_forward: Playing %s.", song.Title))

Outer:
	for {
		err = <-conn.Done
		done, _ := conn.StreamingSession.Finished()
		switch {
		case err.Error() == "stop":
			conn.YoutubeCleanup()
			ctx.Reply(":stop_button: Stopped music. Queue cleared.")
			return
		case err.Error() == "skip":
			conn.Queue.Remove(conn.Queue.Front())
			ctx.Reply(":arrow_forward: Skipped current track.")
			break Outer
		case !done && err != io.EOF:
			conn.YoutubeCleanup()
			log.Printf("player: song failed to play %s", song.Title)
			ctx.Reply(":information_source: Sorry, the music stream broke.")
			return
		case done && err == io.EOF:
			conn.Queue.Remove(conn.Queue.Front())
			break Outer
		}
	}

	go Play(ctx, conn, vc)
}
