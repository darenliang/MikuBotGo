package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/anilist"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"strconv"
	"strings"
)

// Manga command
func Manga(ctx *exrouter.Context) {
	mangaName := strings.TrimSpace(ctx.Args.After(1))

	if mangaName == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Query not specified")
		return
	}

	response := framework.AniListMangaSearchResponse{}
	err := anilist.Query(framework.AnilistMangaSearchQuery(mangaName), &response)

	manga := response.Data.Media

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Manga not found.")
		log.Printf("manga: not found: %s", mangaName)
		return
	}

	volumes := ""
	if manga.Volumes == 0 {
		volumes = "Unknown"
	} else {
		volumes = strconv.Itoa(manga.Volumes)
	}

	published := ""
	if manga.StartDate.Day == 0 {
		published = "Unknown"
	} else {
		published += framework.ParseDate(manga.StartDate.Year, manga.StartDate.Month, manga.StartDate.Day)
		if manga.EndDate.Day != 0 {
			published += " to "
			published += framework.ParseDate(manga.EndDate.Year, manga.EndDate.Month, manga.EndDate.Day)
		}
	}

	genres := "Unknown"
	if len(manga.Genres) != 0 {
		genres = strings.Join(manga.Genres, ", ")
	}

	staffList := make([]string, 0)
	for _, studio := range manga.Staff.Edges {
		staffList = append(staffList, studio.Node.Name.Full)
	}

	staff := "Unknown"
	if len(staffList) != 0 {
		staff = strings.Join(staffList, ", ")
	}

	score := "Unknown"
	if manga.AverageScore != 0 {
		score = strconv.Itoa(manga.AverageScore) + "/100"
	}

	var color uint64
	if manga.CoverImage.Color != "" {
		color, err = strconv.ParseUint(manga.CoverImage.Color[1:], 16, 64)
		if err != nil {
			color = config.EmbedColor
		}
	} else {
		color = config.EmbedColor
	}

	malLink := "Unknown"
	if manga.IDMal != 0 {
		malLink = fmt.Sprintf("https://myanimelist.net/manga/%d", manga.IDMal)
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       int(color),
		Description: strings.ReplaceAll(manga.Description, "<br>", ""),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Type",
				Value:  strings.ReplaceAll(manga.Format, "_", " "),
				Inline: true,
			},
			{
				Name:   "Volumes",
				Value:  volumes,
				Inline: true,
			},
			{
				Name:   "Status",
				Value:  strings.Title(strings.ToLower(manga.Status)),
				Inline: true,
			},
			{
				Name:   "Published",
				Value:  published,
				Inline: false,
			},
			{
				Name:   "Genres",
				Value:  genres,
				Inline: false,
			},
			{
				Name:   "Staff",
				Value:  staff,
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  strings.Title(strings.ToLower(strings.ReplaceAll(manga.Source, "_", " "))),
				Inline: true,
			},
			{
				Name:   "Score",
				Value:  score,
				Inline: true,
			},
			{
				Name:   "MAL Link",
				Value:  malLink,
				Inline: false,
			},
		},
		Title: manga.Title.UserPreferred,
		URL:   manga.SiteURL,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: manga.CoverImage.ExtraLarge,
		},
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
