package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/jozsefsallai/gophersauce"
	"net/url"
	"os"
	"strings"
)

var SauceNaoToken string

func init() {
	SauceNaoToken = os.Getenv("SAUCENAO_TOKEN")
}

func option(value string) string {
	if value == "" {
		return "Unknown"
	}
	return value
}

// Sauce command
func Sauce(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	URL := ""

	if len(query) != 0 {
		_, err := url.ParseRequestURI(query)
		if err != nil {
			ctx.Reply(":thinking: The URL you have provided is not valid.")
			return
		}
		URL = query
	}

	if len(ctx.Msg.Attachments) == 0 && URL == "" {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%ssauce <artwork or screenshot (attachment or url)>`", prefix))
		return
	}

	if URL == "" {
		URL = ctx.Msg.Attachments[0].URL
	}

	client, err := gophersauce.NewClient(&gophersauce.Settings{
		MaxResults: 1,
		APIKey:     SauceNaoToken,
	})

	if err != nil {
		ctx.Reply(":cry: An error has occurred.")
		return
	}

	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	resp, err := client.FromURL(URL)

	if err != nil {
		ctx.Reply(":warning: There's an issue with your image.")
		return
	}

	firstResult := resp.First()

	if len(firstResult.Data.ExternalURLs) == 0 {
		ctx.Reply(":cry: We can't find the sauce.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Title:  fmt.Sprintf("Here's the sauce %s#%s", ctx.Msg.Author.Username, ctx.Msg.Author.Discriminator),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: firstResult.Header.Thumbnail,
		},
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Similarity",
		Value:  fmt.Sprintf("%s%%", firstResult.Header.Similarity),
		Inline: false,
	})

	if firstResult.IsPixiv() {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  option(firstResult.Data.Title),
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  option(firstResult.Data.MemberName),
				Inline: true,
			},
		}...)
	} else if firstResult.IsAniDB() {
		year := firstResult.Data.Year
		splitYear := strings.Split(firstResult.Data.Year, "-")
		if len(splitYear) == 2 && splitYear[0] == splitYear[1] {
			year = splitYear[0]
		}
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Anime Name",
				Value:  option(firstResult.Data.Source),
				Inline: true,
			},
			{
				Name:   "Episode Number",
				Value:  option(firstResult.Data.Part),
				Inline: true,
			},
			{
				Name:   "Year Aired",
				Value:  option(year),
				Inline: true,
			},
		}...)
	} else if firstResult.IsIMDb() {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  option(firstResult.Data.Source),
				Inline: true,
			},
			{
				Name:   "Episode Number",
				Value:  option(firstResult.Data.Part),
				Inline: true,
			},
			{
				Name:   "Year Aired",
				Value:  option(firstResult.Data.Year),
				Inline: true,
			},
		}...)
	} else if firstResult.IsDeviantArt() {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  option(firstResult.Data.Title),
				Inline: true,
			},
			{
				Name:   "Author",
				Value:  option(firstResult.Data.AuthorName),
				Inline: true,
			},
		}...)
	} else if firstResult.IsBcy() {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Title",
				Value:  option(firstResult.Data.Title),
				Inline: true,
			},
			{
				Name:   "Poster",
				Value:  option(firstResult.Data.MemberName),
				Inline: true,
			},
		}...)
	} else if firstResult.IsDanbooru() || firstResult.IsSankaku() {
		embed.Fields = append(embed.Fields, []*discordgo.MessageEmbedField{
			{
				Name:   "Artist",
				Value:  option(fmt.Sprintf("%v", firstResult.Data.Creator)),
				Inline: true,
			},
			{
				Name:   "Source",
				Value:  option(firstResult.Data.Material),
				Inline: true,
			},
			{
				Name:   "Characters",
				Value:  option(firstResult.Data.Characters),
				Inline: true,
			},
		}...)
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Link",
		Value:  firstResult.Data.ExternalURLs[0],
		Inline: false,
	})

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
