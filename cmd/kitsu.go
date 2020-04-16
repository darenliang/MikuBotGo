package cmd

import (
	"encoding/json"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/kitsu"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"net/url"
	"strconv"
	"strings"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

// Anime command
func Anime(ctx *exrouter.Context) {
	animeName := strings.TrimSpace(ctx.Args.After(1))

	if animeName == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Query not specified")
		return
	}

	search, _ := kitsu.GetAnimePage(`anime?filter[text]=` + url.QueryEscape(animeName) + `&page[limit]=1`)

	if len(search.Data) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Anime not found")
		return
	}

	anime := search.Data[0]

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: search.Data[0].Attributes.Synopsis,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Type",
				Value:  anime.Attributes.Subtype,
				Inline: true,
			},
			{
				Name:   "Episodes",
				Value:  strconv.Itoa(anime.Attributes.EpisodeCount),
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  strings.Title(anime.Attributes.Status),
				Inline: true,
			},
			{
				Name:   "Premiered",
				Value:  anime.Attributes.StartDate,
				Inline: true,
			},
			{
				Name:   "Score",
				Value:  anime.Attributes.AverageRating,
				Inline: true,
			},
			{
				Name:   "Rating",
				Value:  anime.Attributes.AgeRating,
				Inline: true,
			},
		},
		Title: anime.Attributes.CanonicalTitle,
		URL:   anime.Link(),
		Image: &discordgo.MessageEmbedImage{
			URL: anime.Attributes.PosterImage.Original,
		},
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
