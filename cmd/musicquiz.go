package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/jikan-go"
	"math/rand"
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
	guess := ctx.Args.After(1)
	guess = strings.TrimSpace(guess)

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
	fileName := fmt.Sprintf("%s.mp4", config.Openings[idx].File)

	response, _ := jikan.Search{Type: "anime", Q: config.Openings[idx].Source}.Get()

	answers := make([]string, 3)
	for i := 0; i < 3; i++ {
		answers = append(answers, response["results"].([]interface{})[i].(map[string]interface{})["title"].(string))
	}

	config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{
		Answers: answers,
		Source:  config.Openings[idx].Source,
	}

	fileNameOut := fmt.Sprintf("%s.mp3", framework.RandomString(16))

	_ = framework.DownloadFile(fmt.Sprintf("./cache/%s", fileName),
		fmt.Sprintf("https://openings.moe/video/%s", fileName))

	cmd := exec.Command("ffmpeg", "-i", "./cache/"+fileName,
		"-vn", "-ab", "128k", "-ar", "44100", "-y", "./cache/"+fileNameOut)

	ch := make(chan error)
	go func() {
		ch <- cmd.Run()
	}()

	select {
	case err := <-ch:
		_ = os.Remove("./cache/" + fileName)
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Failed to convert media file.")
			_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f", ctx.Ses.State.User.ID)
			return
		}
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf(
			"`%smusicquiz answer` to guess anime or `%smusicquiz giveup` to give up.", config.Prefix, config.Prefix))
		f, _ := os.Open("./cache/" + fileNameOut)
		_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fileNameOut, f)
		_ = f.Close()
		_ = os.Remove("./cache/" + fileNameOut)
	case <-time.After(config.Timeout * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "This command timeout.")
	}

	_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f", ctx.Ses.State.User.ID)
}
