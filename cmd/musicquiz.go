package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/anilist"
	"github.com/animenotifier/kitsu"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strconv"
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
		if config.OpeningsMap[ctx.Msg.ChannelID].Name == "" {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "There is no currently active quiz in this channel.")
			return
		} else {
			if guess == "giveup" {
				config.OpeningsMap[ctx.Msg.ChannelID].Embed.Title = "Answer: " + config.OpeningsMap[ctx.Msg.ChannelID].Embed.Title
				config.OpeningsMap[ctx.Msg.ChannelID].Embed.Color = 0xf44336
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, config.OpeningsMap[ctx.Msg.ChannelID].Embed)
				config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningEntry{}
				return
			} else if guess == "hint" {
				response := framework.AniListAnimeSearchResponse{}
				_ = anilist.Query(framework.AnilistAnimeSearchQuery(config.OpeningsMap[ctx.Msg.ChannelID].Name), &response)

				anime := response.Data.Media

				properStudios := make([]string, 0)
				for _, studio := range anime.Studios.Edges {
					if studio.Node.IsAnimationStudio {
						properStudios = append(properStudios, studio.Node.Name)
					}
				}

				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  config.EmbedColor,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name: "Season",
							Value: fmt.Sprintf("%s %d", strings.Title(strings.ToLower(anime.Season)),
								anime.SeasonYear),
							Inline: false,
						},
						{
							Name:   "Genres",
							Value:  strings.Join(anime.Genres, ", "),
							Inline: false,
						},
						{
							Name:   "Studios",
							Value:  strings.Join(properStudios, ", "),
							Inline: false,
						},
					},
					Title: "Hints for musicquiz",
				}
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
				return
			} else {
				response, _ := kitsu.GetAnimePage(`anime?filter[text]=` + url.QueryEscape(guess) + `&page[limit]=3`)
				answers := make([]string, 0)
				for _, val := range response.Data {
					answers = append(answers, val.Attributes.Titles.En)
					answers = append(answers, val.Attributes.Titles.EnJp)
				}
				if framework.GetStringValidation(answers, config.OpeningsMap[ctx.Msg.ChannelID].Name) {
					config.OpeningsMap[ctx.Msg.ChannelID].Embed.Title = "Correct: " + config.OpeningsMap[ctx.Msg.ChannelID].Embed.Title
					config.OpeningsMap[ctx.Msg.ChannelID].Embed.Color = 0x4caf50
					_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, config.OpeningsMap[ctx.Msg.ChannelID].Embed)
					score, attempts := framework.GetDatabaseValue(ctx.Msg.Author.ID)
					if score == 0 && attempts == 0 {
						framework.CreateDatabaseEntry(ctx.Msg.Author.ID, 1, 0)
					} else {
						framework.UpdateDatabaseValue(ctx.Msg.Author.ID, score+1, attempts)
					}
					config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningEntry{}
					return
				} else {
					_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You are incorrect. Please try again.")
					return
				}
			}
		}
	}

	if config.OpeningsMap[ctx.Msg.ChannelID].Name != "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You haven't gave an answer to the previous quiz.")
		return
	}

	score, attempts := framework.GetDatabaseValue(ctx.Msg.Author.ID)
	if score == 0 && attempts == 0 {
		framework.CreateDatabaseEntry(ctx.Msg.Author.ID, 0, 1)
	} else {
		framework.UpdateDatabaseValue(ctx.Msg.Author.ID, score, attempts+1)
	}

	_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f")

	idx := rand.Int() % len(config.Openings)

	response := framework.AniListAnimeSearchResponse{}
	_ = anilist.Query(framework.AnilistAnimeSearchQuery(config.Openings[idx].Name), &response)

	fileName := config.Openings[idx].Songs[rand.Int()%len(config.Openings[idx].Songs)]

	// This is whack-----------------------------------------------
	r, _ := httpClient.Get(fmt.Sprintf("https://api.jikan.moe/v3/anime/%d", response.Data.Media.IDMal))
	anime := framework.MALAnime{}
	_ = json.NewDecoder(r.Body).Decode(&anime)

	re := regexp.MustCompile("-(\\D+)(\\d+)v*\\d*.webm$")
	match := re.FindStringSubmatch(fileName)

	songIdx, _ := strconv.Atoi(match[2])

	song := ""
	if strings.HasPrefix(match[1], "OP") {
		if len(anime.OpeningThemes) <= songIdx-1 {
			song = "Cannot find song info"
		} else {
			song = anime.OpeningThemes[songIdx-1]
		}
	} else {
		if len(anime.EndingThemes) <= songIdx-1 {
			song = "Cannot find song info"
		} else {
			song = anime.EndingThemes[songIdx-1]
		}
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Song",
				Value:  song,
				Inline: true,
			},
		},
		Title: response.Data.Media.Title.UserPreferred,
		URL:   fmt.Sprintf("https://myanimelist.net/anime/%d", response.Data.Media.IDMal),
	}

	// This is whack-----------------------------------------------

	config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningEntry{
		Name:  response.Data.Media.Title.UserPreferred,
		Embed: embed,
	}

	fileNameOut := framework.RandomString(16)

	cmd := exec.Command("youtube-dl", "--extract-audio", "--audio-format", "mp3", "--output",
		"./cache/"+fileNameOut+".webm", "--external-downloader", "aria2c", "--external-downloader-args",
		`-x 5 -s 5 -k 1M`, fileName)

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
			"`%smusicquiz <guess>` to guess anime, `%smusicquiz hint` to get hints or `%smusicquiz giveup` to give up.", config.Prefix, config.Prefix, config.Prefix))
		f, err := os.Open("./cache/" + fileNameOut + ".mp3")
		_, err = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fileNameOut+".mp3", f)
		_ = f.Close()
		_ = os.Remove("./cache/" + fileNameOut + ".mp3")
	case <-time.After(config.Timeout * 3 * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "This command timed out.")
		config.OpeningsMap[ctx.Msg.ChannelID] = config.OpeningEntry{}
	}

	_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, "\xe2\x8f\xb2\xef\xb8\x8f", ctx.Ses.State.User.ID)
}
