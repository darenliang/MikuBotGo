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
	"io/ioutil"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	OpeningsData Openings
	OpeningsMap  = sync.Map{}
)

type OpeningEntry struct {
	Name  string
	Embed *discordgo.MessageEmbed
}

type Openings []struct {
	Name  string `json:"name"`
	Songs []struct {
		Songname string `json:"songname"`
		URL      string `json:"url"`
	} `json:"songs"`
}

func init() {
	// Generate random seed
	rand.Seed(time.Now().UnixNano())

	// Setup openings
	OpeningsData = GetOpenings()
}

// Return openings
func GetOpenings() Openings {
	file, _ := ioutil.ReadFile("data/dataset_filtered.json")
	tmp := Openings{}
	_ = json.Unmarshal(file, &tmp)
	return tmp
}

// MusicQuiz command
func MusicQuiz(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	guess := strings.TrimSpace(ctx.Args.After(1))

	if len(guess) != 0 {
		entryInterface, ok := OpeningsMap.Load(ctx.Msg.ChannelID)
		if !ok {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "There is no currently active quiz in this channel.")
			return
		} else {
			entry := entryInterface.(OpeningEntry)
			if guess == "giveup" {
				entry.Embed.Title = "Answer: " + entry.Embed.Title
				entry.Embed.Color = 0xf44336
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, entry.Embed)
				OpeningsMap.Delete(ctx.Msg.ChannelID)
				return
			} else if guess == "hint" {
				response := framework.AniListAnimeSearchResponse{}
				_ = anilist.Query(framework.AnilistAnimeSearchQuery(entry.Name), &response)

				anime := response.Data.Media

				properStudios := make([]string, 0)
				for _, studio := range anime.Studios.Edges {
					if studio.Node.IsAnimationStudio {
						properStudios = append(properStudios, studio.Node.Name)
					}
				}

				studios := "Unknown"
				if len(properStudios) != 0 {
					studios = strings.Join(properStudios, ", ")
				}

				genres := "Unknown"
				if len(anime.Genres) != 0 {
					genres = strings.Join(anime.Genres, ", ")
				}

				embed := &discordgo.MessageEmbed{
					Author: &discordgo.MessageEmbedAuthor{},
					Color:  config.EmbedColor,
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Type",
							Value:  strings.ReplaceAll(anime.Format, "_", " "),
							Inline: true,
						},
						{
							Name: "Season",
							Value: fmt.Sprintf("%s %d", strings.Title(strings.ToLower(anime.Season)),
								anime.SeasonYear),
							Inline: true,
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
					Title: "Hints for Music Quiz",
				}
				_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
				return
			} else {
				response, _ := kitsu.GetAnimePage(`anime?filter[text]=` + url.QueryEscape(guess) + `&page[limit]=5`)
				answers := make([]string, 0)
				for _, val := range response.Data {
					answers = append(answers, val.Attributes.Titles.En, val.Attributes.Titles.EnJp)
				}
				if framework.GetStringValidation(answers, entry.Name) {
					entry.Embed.Title = "Correct: " + entry.Embed.Title
					entry.Embed.Color = 0x4caf50
					_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, entry.Embed)
					score, attempts := framework.MQDB.GetScore(ctx.Msg.Author.ID)
					if score == 0 && attempts == 0 {
						framework.MQDB.CreateScore(ctx.Msg.Author.ID, 1, 0)
					} else {
						framework.MQDB.UpdateScore(ctx.Msg.Author.ID, score+1, attempts)
					}
					OpeningsMap.Delete(ctx.Msg.ChannelID)
					return
				} else {
					_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You are incorrect. Please try again.")
					return
				}
			}
		}
	}

	_, ok := OpeningsMap.Load(ctx.Msg.ChannelID)
	if ok {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You haven't gave an answer to the previous quiz.")
		return
	}

	score, attempts := framework.MQDB.GetScore(ctx.Msg.Author.ID)
	if score == 0 && attempts == 0 {
		framework.MQDB.CreateScore(ctx.Msg.Author.ID, 0, 1)
	} else {
		framework.MQDB.UpdateScore(ctx.Msg.Author.ID, score, attempts+1)
	}

	_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	idx := rand.Int() % len(OpeningsData)

	response := framework.AniListAnimeSearchResponse{}
	_ = anilist.Query(framework.AnilistAnimeSearchQuery(OpeningsData[idx].Name), &response)

	song := OpeningsData[idx].Songs[rand.Int()%len(OpeningsData[idx].Songs)]

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Song",
				Value:  song.Songname,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: response.Data.Media.CoverImage.ExtraLarge,
		},
		Title: response.Data.Media.Title.UserPreferred,
		URL:   fmt.Sprintf("https://myanimelist.net/anime/%d", response.Data.Media.IDMal),
	}

	OpeningsMap.Store(ctx.Msg.ChannelID, OpeningEntry{
		Name:  response.Data.Media.Title.UserPreferred,
		Embed: embed,
	})

	fileNameOut := framework.RandomString(16)

	resp, err := framework.HttpClient.Get(fmt.Sprintf(
		"https://gitlab.com/darenliang/mq/-/raw/master/data/%s", song.URL))

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Failed to get media file.")
		_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
		log.Printf("musicquiz: file failed to convert: %s", song.URL)
		OpeningsMap.Delete(ctx.Msg.ChannelID)
		return
	}

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf(
		"`%smusicquiz <guess>` to guess anime, `%smusicquiz hint` to get hints or `%smusicquiz giveup` to give up.", prefix, prefix, prefix))
	_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fileNameOut+".mp3", resp.Body)
	_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
}
