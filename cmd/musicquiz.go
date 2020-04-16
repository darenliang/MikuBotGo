package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
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
			} else if guess == "hint" {
				anime, _ := jikan.Anime{ID: config.OpeningsMap[ctx.Msg.ChannelID].Id}.Get()

				var premiered string
				if anime["premiered"] == nil {
					premiered = "Not available"
				} else {
					premiered = anime["premiered"].(string)
				}

				genres := ""
				if len(anime["genres"].([]interface{})) != 0 {
					end := len(anime["genres"].([]interface{}))
					for idx, genre := range anime["genres"].([]interface{}) {
						genres += genre.(map[string]interface{})["name"].(string)
						if idx != end-1 {
							genres += ", "
						}
					}
				} else {
					genres = "Not available"
				}

				studios := ""
				if len(anime["studios"].([]interface{})) != 0 {
					end := len(anime["studios"].([]interface{}))
					for idx, studio := range anime["studios"].([]interface{}) {
						studios += studio.(map[string]interface{})["name"].(string)
						if idx != end-1 {
							studios += ", "
						}
					}
				} else {
					studios = "Not available"
				}

				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  config.EmbedColor,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Premiered",
							Value:  premiered,
							Inline: false,
						},
						{
							Name:   "Genres",
							Value:  genres,
							Inline: false,
						},
						{
							Name:   "Studios",
							Value:  studios,
							Inline: false,
						},
					},
					Timestamp: time.Now().Format(time.RFC3339),
					Title:     "Hints for musicquiz",
				}
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
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

	idx := rand.Int() % len(config.Openings)
	fileName := fmt.Sprintf("%s.webm", config.Openings[idx].File)

	response, _ := jikan.Search{Type: "anime", Q: config.Openings[idx].Source}.Get()

	answers := make([]string, 3)
	for i := 0; i < 3; i++ {
		answers = append(answers, response["results"].([]interface{})[i].(map[string]interface{})["title"].(string))
	}

	config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{
		Id:      int(response["results"].([]interface{})[0].(map[string]interface{})["mal_id"].(float64)),
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
			return
		}
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf(
			"`%smusicquiz answer` to guess anime, `%smusicquiz hint` to get hints or `%smusicquiz giveup` to give up.", config.Prefix, config.Prefix, config.Prefix))
		f, err := os.Open("./cache/" + fileNameOut + ".mp3")
		_, err = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fileNameOut+".mp3", f)
		_ = f.Close()
		_ = os.Remove("./cache/" + fileNameOut + ".mp3")
	case <-time.After(config.Timeout * 3 * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "This command timed out.")
		config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningsEntry{}
	}
}
