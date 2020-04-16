package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/kitsu"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func init() {
	// Generate random seed
	rand.Seed(time.Now().UnixNano())
}

// MusicQuiz command
func MusicQuiz(ctx *exrouter.Context) {
	guess := strings.TrimSpace(ctx.Args.After(1))

	if len(guess) != 0 {
		if config.OpeningsMap[ctx.Msg.ChannelID].Source == "" {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "There is no currently active quiz in this channel.")
			return
		} else {
			if guess == "giveup" {
				_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "The answer is "+config.OpeningsMap[ctx.Msg.ChannelID].Source)
				config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{}
				return
			} else if framework.GetStringValidation(config.OpeningsMap[ctx.Msg.ChannelID].Answers, guess) {
				_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You are correct! The answer is "+config.OpeningsMap[ctx.Msg.ChannelID].Source)
				config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{}
				score := framework.GetDatabaseValue(ctx.Msg.Author.ID)
				if score == 0 {
					framework.CreateDatabaseEntry(ctx.Msg.Author.ID, 1)
				} else {
					framework.UpdateDatabaseValue(ctx.Msg.Author.ID, score+1)
				}
				return
			} else {
				_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You are incorrect. Please try again.")
				return
			}
		}
	}

	if config.OpeningsMap[ctx.Msg.ChannelID].Source != "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You haven't gave an answer to the previous quiz.")
		return
	}

	_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f")

	idx := rand.Int() % len(config.Openings)
	fileName := fmt.Sprintf("%s.webm", config.Openings[idx].File)

	response, _ := kitsu.GetAnimePage(`anime?filter[text]=` + url.QueryEscape(config.Openings[idx].Source) + `&page[limit]=3`)

	answers := make([]string, 0)
	for _, val := range response.Data {
		answers = append(answers, val.Attributes.CanonicalTitle)
	}

	config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{
		Answers: answers,
		Source:  config.Openings[idx].Source,
	}

	fileNameOut := framework.RandomString(16)

	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "--output",
		"./cache/"+fileNameOut+".webm", "--external-downloader", "aria2c", "--external-downloader-args",
		`-x 5 -s 5 -k 1M`, "https://openings.moe/video/"+fileName)

	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	select {
	case err := <-ch:
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Failed to convert media file.")
			_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f", ctx.Ses.State.User.ID)
			return
		}
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf(
			"`%smusicquiz answer` to guess anime or `%smusicquiz giveup` to give up.", config.Prefix, config.Prefix))
		f, err := os.Open("./cache/" + fileNameOut + ".mp3")
		_, err = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fileNameOut+".mp3", f)
		_ = f.Close()
		_ = os.Remove("./cache/" + fileNameOut + ".mp3")
	case <-time.After(config.Timeout * 3 * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "This command timed out.")
		config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{}
	}

	_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f", ctx.Ses.State.User.ID)
}
