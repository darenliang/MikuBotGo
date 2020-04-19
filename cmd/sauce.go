package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/anilist"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type TraceData struct {
	Docs []struct {
		AnilistID int         `json:"anilist_id"`
		Episode   interface{} `json:"episode,string"`
		MalID     int         `json:"mal_id"`
	} `json:"docs"`
}

var httpClient = &http.Client{Timeout: config.Timeout * time.Second}

func getJson(url string, target interface{}) error {
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(target)
	if err != nil {
		return err
	}

	err = r.Body.Close()
	if err != nil {
		return err
	}

	return nil
}

// Sauce command
func Sauce(ctx *exrouter.Context) {
	query := strings.TrimSpace(ctx.Args.After(1))

	URL := ""

	if len(query) != 0 {
		_, err := url.ParseRequestURI(query)
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "The URL you have provided is not valid.")
			return
		}
		URL = query
	}

	if len(ctx.Msg.Attachments) == 0 && URL == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You did not attach any links or images.")
		return
	}

	if URL == "" {
		URL = ctx.Msg.Attachments[0].URL
	}

	trace := TraceData{}
	err := getJson("https://trace.moe/api/search?url="+URL, &trace)
	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	if len(trace.Docs) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "We can't find the sauce for you.")
		return
	}

	response := framework.TraceSearchResult{}
	_ = anilist.Query(framework.AnilistAnimeIDQuery(trace.Docs[0].AnilistID), &response)
	anime := response.Data.Media

	var color uint64
	if anime.CoverImage.Color != "" {
		color, err = strconv.ParseUint(anime.CoverImage.Color[1:], 16, 64)
		if err != nil {
			color = config.EmbedColor
		}
	} else {
		color = config.EmbedColor
	}

	episode := fmt.Sprint(trace.Docs[0].Episode)
	if episode == "" {
		episode = "Unknown"
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  int(color),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Anime",
				Value:  anime.Title.UserPreferred,
				Inline: true,
			},
			{
				Name:   "Episode Number",
				Value:  episode,
				Inline: true,
			},
			{
				Name:   "MAL Link",
				Value:  fmt.Sprintf("https://myanimelist.net/anime/%d", trace.Docs[0].MalID),
				Inline: false,
			},
		},
		Title: fmt.Sprintf("Here's the sauce %s#%s", ctx.Msg.Author.Username, ctx.Msg.Author.Discriminator),
		Image: &discordgo.MessageEmbedImage{
			URL: anime.CoverImage.ExtraLarge,
		},
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
