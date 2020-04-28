// Code taken from ducc/GoMusicBot
package cmd

import (
	"bytes"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"math/rand"
	"strconv"
	"strings"
)

const (
	songFormat    = "\n`%03d` %s"
	currentFormat = "__Current song__\n%s\n"
	invalidPage   = "Invalid page `%d`. Min: `1`, max: `%d`"
)

// Add music command
func AddMusic(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))
	if len(query) == 0 {
		_, _ = ctx.Reply("Please provide a query.")
		return
	}
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	msg, _ := ctx.Reply("Adding songs to queue...")
	for idx, arg := range ctx.Args {
		if idx == 0 {
			continue
		}
		t, inp, err := music.Youtube{}.Get(arg)

		if err != nil {
			_, _ = ctx.Reply("An error occurred.")
			return
		}

		switch t {
		case music.ERRORTYPE:
			_, _ = ctx.Reply("An error occurred.")
			return
		case music.VIDEOTYPE:
			{
				video, err := music.Youtube{}.Video(*inp)
				if err != nil {
					_, _ = ctx.Reply("Cannot find music.")
					return
				}
				song := music.NewSong(video.Media, video.Title, arg)
				musicSession.Queue.Add(*song)
				_, _ = ctx.Ses.ChannelMessageEdit(ctx.Msg.ChannelID, msg.ID, "Added `"+song.Title+"` to the song queue."+
					fmt.Sprintf(" Use `%splay` to start playing the songs. To see the song queue, use `%squeue`.", prefix, prefix))
				break
			}
		case music.PLAYLISTTYPE:
			{
				videos, err := music.Youtube{}.Playlist(*inp)
				if err != nil {
					_, _ = ctx.Reply("An error occurred.")
					return
				}
				for _, v := range *videos {
					id := v.Id
					_, i, err := music.Youtube{}.Get(id)
					if err != nil {
						_, _ = ctx.Reply("An error occurred.")
						continue
					}
					video, err := music.Youtube{}.Video(*i)
					if err != nil {
						_, _ = ctx.Reply("Cannot find music.")
						return
					}
					song := music.NewSong(video.Media, video.Title, arg)
					musicSession.Queue.Add(*song)
				}
				_, _ = ctx.Reply(fmt.Sprintf("Finished adding songs to the playlist. Use `%splay` to start playing the songs. To see the song queue, use `%squeue`.", prefix, prefix))
				break
			}
		}
	}
}

// Clear music command
func ClearCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	if !musicSession.Queue.HasNext() {
		_, _ = ctx.Reply("Queue is already empty.")
		return
	}
	musicSession.Queue.Clear()
	_, _ = ctx.Reply("Cleared the song queue.")
}

// Current music command
func CurrentCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	current := musicSession.Queue.Current()
	if current == nil {
		_, _ = ctx.Reply(fmt.Sprintf("The song queue is empty. Add a song with `%sadd`.", prefix))
		return
	}
	_, _ = ctx.Reply("Currently playing `" + current.Title + "`.")
}

// Join music command
func JoinCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	if music.MusicSessions.GetByGuild(ctx.Msg.GuildID) != nil {
		_, _ = ctx.Reply(fmt.Sprintf("Already connected. Use `%sleave` for the bot to disconnect.", prefix))
		return
	}
	guild, err := ctx.Guild(ctx.Msg.GuildID)
	if err != nil {
		_, _ = ctx.Reply("An error occurred.")
		return
	}

	var voiceChannel *discordgo.Channel
	for _, state := range guild.VoiceStates {
		if state.UserID == ctx.Msg.Author.ID {
			channel, _ := ctx.Ses.State.Channel(state.ChannelID)
			voiceChannel = channel
			break
		}
	}

	if voiceChannel == nil {
		_, _ = ctx.Reply("You must be in a voice channel to use the bot.")
		return
	}

	sess, err := music.MusicSessions.Join(ctx.Ses, ctx.Msg.GuildID, voiceChannel.ID, music.JoinProperties{
		Muted:    false,
		Deafened: true,
	})

	if err != nil {
		_, _ = ctx.Reply("An error occurred.")
		return
	}

	_, _ = ctx.Reply("Joined <#" + sess.ChannelId + ">.")

	// Handle timeout
	// go music.HandleMusicTimeout(sess, func(msg string) {
	// 	_, _ = ctx.Reply(msg)
	// })
}

// Leave music command
func LeaveCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	music.MusicSessions.Leave(ctx.Ses, *musicSession)
	_, _ = ctx.Reply("Left <#" + musicSession.ChannelId + ">.")
}

// Pause music command
func PauseCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	queue := musicSession.Queue
	if !queue.HasNext() {
		_, _ = ctx.Reply(fmt.Sprintf("Queue is empty. Add songs with `%sadd`.", prefix))
		return
	}
	queue.Pause()

	_, _ = ctx.Reply(fmt.Sprintf("The queue has paused and will stop playing after this song. To resume the queue, use `%splay`.", prefix))

	// Handle timeout
	// go music.HandleMusicTimeout(musicSession, func(msg string) {
	// 	_, _ = ctx.Reply(msg)
	// })
}

// Play music command
func PlayCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	queue := musicSession.Queue
	if !queue.HasNext() {
		_, _ = ctx.Reply(fmt.Sprintf("Queue is empty. Add songs with `%sadd`.", prefix))
		return
	}
	go queue.Start(musicSession, func(msg string) {
		_, _ = ctx.Reply(msg)
	})
}

// Queue music command
func QueueCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	queue := musicSession.Queue
	q := queue.Get()
	if len(q) == 0 && queue.Current() == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Song queue is empty. Add a song with `%sadd`.", prefix))
		return
	}
	buff := bytes.Buffer{}
	if queue.Current() != nil {
		buff.WriteString(fmt.Sprintf(currentFormat, queue.Current().Title))
	}
	queueLength := len(q)
	if len(ctx.Args) == 0 {
		var resp string
		if queueLength > 20 {
			resp = display(q[:20], buff, 2, 0)
		} else {
			resp = display(q[:queueLength], buff, 2, 0)
		}
		_, _ = ctx.Reply(resp)
		return
	}
	page, err := strconv.Atoi(ctx.Args[0])
	if err != nil {
		_, _ = ctx.Reply("Invalid page `" + ctx.Args[0] + fmt.Sprintf("`. Usage: `%squeue <page>`", prefix))
		return
	}
	pages := queueLength / 20
	if page < 1 || page > (pages+1) {
		_, _ = ctx.Reply(fmt.Sprintf(invalidPage, page, pages+1))
		return
	}
	var lowerBound int
	if page == 1 {
		lowerBound = 0
	} else {
		lowerBound = (page - 1) * 20
	}
	upperBound := page * 20
	if upperBound > queueLength {
		upperBound = queueLength
	}
	slice := q[lowerBound:upperBound]
	_, _ = ctx.Reply(display(slice, buff, page+1, lowerBound))
}

func display(queue []music.Song, buff bytes.Buffer, page, start int) string {
	for index, song := range queue {
		buff.WriteString(fmt.Sprintf(songFormat, start+index+1, song.Title))
	}
	buff.WriteString(fmt.Sprintf("\n\nView the next page: `queue %d`", page))
	return buff.String()
}

// Shuffle music command
func ShuffleCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	queue := musicSession.Queue
	if !queue.HasNext() {
		_, _ = ctx.Reply(fmt.Sprintf("Queue is empty. Add songs with `%sadd`.", prefix))
		return
	}
	dest := shuffleLoop(queue.Get(), 3)
	queue.Set(dest)
	_, _ = ctx.Reply("Shuffled the song queue.")
}

func shuffleLoop(list []music.Song, i int) []music.Song {
	for x := 0; x < i; x++ {
		list = shuffle(list)
	}
	return list
}

func shuffle(list []music.Song) []music.Song {
	dest := make([]music.Song, len(list))
	perm := rand.Perm(len(list))
	for i, v := range perm {
		dest[v] = list[i]
	}
	return dest
}

// Skip music command
func SkipCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	musicSession.Stop()
	_, _ = ctx.Reply("Skipped song.")
}

// Stop music command
func StopCommand(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	musicSession := music.MusicSessions.GetByGuild(ctx.Msg.GuildID)
	if musicSession == nil {
		_, _ = ctx.Reply(fmt.Sprintf("Not in a voice channel. To make the bot join one, use `%sjoin`.", prefix))
		return
	}
	if musicSession.Queue.HasNext() {
		musicSession.Queue.Clear()
	}
	musicSession.Stop()

	// Handle timeout
	// go music.HandleMusicTimeout(musicSession, func(msg string) {
	// 	_, _ = ctx.Reply(msg)
	// })
}
