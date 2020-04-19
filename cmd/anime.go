package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/anilist"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"strconv"
	"strings"
)

// Anime command
func Anime(ctx *exrouter.Context) {
	animeName := strings.TrimSpace(ctx.Args.After(1))

	if animeName == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Query not specified")
		return
	}

	response := framework.AniListAnimeSearchResponse{}
	err := anilist.Query(framework.AnilistAnimeSearchQuery(animeName), &response)

	anime := response.Data.Media

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Anime not found.")
	}

	episodes := ""
	if anime.Episodes == 0 {
		episodes = "Unknown"
	} else {
		episodes = strconv.Itoa(anime.Episodes)
	}

	aired := ""
	if anime.StartDate.Day == 0 {
		aired = "Unknown"
	} else {
		aired += framework.ParseDate(anime.StartDate.Year, anime.StartDate.Month, anime.StartDate.Day)
		if anime.EndDate.Day != 0 {
			aired += " to "
			aired += framework.ParseDate(anime.EndDate.Year, anime.EndDate.Month, anime.EndDate.Day)
		}
	}

	genres := "Unknown"
	if len(anime.Genres) != 0 {
		genres = strings.Join(anime.Genres, ", ")
	}

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

	score := "Unknown"
	if anime.AverageScore != 0 {
		score = strconv.Itoa(anime.AverageScore) + "/100"
	}

	var color uint64
	if anime.CoverImage.Color != "" {
		color, err = strconv.ParseUint(anime.CoverImage.Color[1:], 16, 64)
		if err != nil {
			color = config.EmbedColor
		}
	} else {
		color = config.EmbedColor
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       int(color),
		Description: strings.ReplaceAll(anime.Description, "<br>", ""),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Type",
				Value:  strings.ReplaceAll(anime.Format, "_", " "),
				Inline: true,
			},
			{
				Name:   "Episodes",
				Value:  episodes,
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  strings.Title(strings.ToLower(anime.Status)),
				Inline: true,
			},
			{
				Name:   "Aired",
				Value:  aired,
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
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  strings.Title(strings.ToLower(strings.ReplaceAll(anime.Source, "_", " "))),
				Inline: true,
			},
			{
				Name:   "Score",
				Value:  score,
				Inline: true,
			},
			{
				Name:   "MAL Link",
				Value:  fmt.Sprintf("https://myanimelist.net/anime/%d", anime.IDMal),
				Inline: false,
			},
		},
		Title: anime.Title.UserPreferred,
		URL:   anime.SiteURL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.CoverImage.ExtraLarge,
		},
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
