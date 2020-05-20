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

	if len(resp.First().Data.ExternalURLs) == 0 {
		ctx.Reply(":cry: We can't find the sauce.")
		return
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Similarity",
				Value:  fmt.Sprintf("%s%%", resp.First().Header.Similarity),
				Inline: true,
			},
			{
				Name:   "Sauce",
				Value:  resp.First().Data.ExternalURLs[0],
				Inline: true,
			},
		},
		Title: fmt.Sprintf("Here's the sauce %s#%s", ctx.Msg.Author.Username, ctx.Msg.Author.Discriminator),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: resp.First().Header.Thumbnail,
		},
	}

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
