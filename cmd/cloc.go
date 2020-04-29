package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"net/url"
	"strings"
)

type ClocResponse []struct {
	Language    string `json:"language"`
	Files       string `json:"files"`
	Lines       string `json:"lines"`
	Blanks      string `json:"blanks"`
	Comments    string `json:"comments"`
	LinesOfCode string `json:"linesOfCode"`
}

// Cloc command
func Cloc(ctx *exrouter.Context) {
	repo := strings.TrimSpace(ctx.Args.After(1))

	if len(repo) == 0 {
		repo = "darenliang/MikuBotGo"
	}

	cloc := ClocResponse{}
	err := framework.UrlToStruct("https://api.codetabs.com/v1/loc?github="+url.QueryEscape(repo), &cloc)

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot find the GitHub repo you are looking for.")
		log.Printf("cloc: repo not found: %s", repo)
		return
	}

	if len(cloc) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "No programming languages detected.")
		log.Printf("cloc: repo not found: %s", repo)
		return
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: "Only the top 15 languages are listed.",
		Fields:      make([]*discordgo.MessageEmbedField, 0),
		Title:       "How many lines of code in " + repo + "?",
		URL:         "https://github.com/" + repo,
	}

	for idx, val := range cloc {
		if idx == 15 || idx == len(cloc)-1 {
			break
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   val.Language,
			Value:  val.LinesOfCode,
			Inline: true,
		})
	}

	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   cloc[len(cloc)-1].Language,
		Value:  cloc[len(cloc)-1].LinesOfCode,
		Inline: false,
	})

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
