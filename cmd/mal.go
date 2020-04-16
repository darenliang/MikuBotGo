package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/jikan-go"
	"strconv"
	"sync"
	"time"
)

var emojis = []string{
	"\x31\xef\xb8\x8f\xe2\x83\xa3",
	"\x32\xef\xb8\x8f\xe2\x83\xa3",
	"\x33\xef\xb8\x8f\xe2\x83\xa3",
	"\x34\xef\xb8\x8f\xe2\x83\xa3",
	"\x35\xef\xb8\x8f\xe2\x83\xa3",
}

// Anime command
func Anime(ctx *exrouter.Context) {
	var (
		lock     sync.RWMutex
		embedMsg *discordgo.Message
		callback = make(chan struct{})
		idx      int
	)

	animeName := ctx.Args.After(1)
	search, _ := jikan.Search{Type: "anime", Q: animeName}.Get()
	results := search["results"].([]interface{})

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Description: fmt.Sprintf(
			":one: %s\n"+
				":two: %s\n"+
				":three: %s\n"+
				":four: %s\n"+
				":five: %s",
			results[0].(map[string]interface{})["title"].(string),
			results[1].(map[string]interface{})["title"].(string),
			results[2].(map[string]interface{})["title"].(string),
			results[3].(map[string]interface{})["title"].(string),
			results[4].(map[string]interface{})["title"].(string),
		),
		Timestamp: time.Now().Format(time.RFC3339),
		Title:     fmt.Sprintf("Search results for %s", animeName),
	}

	embedMsg, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	for i := 0; i < 5; i++ {
		_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, embedMsg.ID, emojis[i])
	}

	defer ctx.Ses.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		lock.RLock()
		defer lock.RUnlock()

		idx = framework.Index(reaction.Emoji.Name, emojis)

		// TODO: Use for debug
		// fmt.Println(hex.EncodeToString([]byte(reaction.Emoji.Name)))

		if reaction.MessageID == embedMsg.ID && reaction.UserID != ctx.Ses.State.User.ID && idx != -1 {
			close(callback)
		}
	})()

	select {
	case <-callback:
		malID := int(results[idx].(map[string]interface{})["mal_id"].(float64))
		anime, _ := jikan.Anime{ID: malID}.Get()
		title := anime["title"].(string)

		var episodes string
		if anime["episodes"] == nil {
			episodes = "Not available"
		} else {
			episodes = strconv.Itoa(int(anime["episodes"].(float64)))
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

		var score string
		if anime["score"] == nil {
			score = "Not available"
		} else {
			score = fmt.Sprintf("%.2f", anime["score"].(float64))
		}

		var rank string
		if anime["rank"] == nil {
			rank = "Not available"
		} else {
			rank = strconv.Itoa(int(anime["rank"].(float64)))
		}

		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       config.EmbedColor,
			Description: anime["synopsis"].(string),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Type",
					Value:  anime["type"].(string),
					Inline: true,
				},
				{
					Name:   "Episodes",
					Value:  episodes,
					Inline: true,
				},
				{
					Name:   "Status",
					Value:  anime["status"].(string),
					Inline: true,
				},
				{
					Name:   "Aired",
					Value:  anime["aired"].(map[string]interface{})["string"].(string),
					Inline: false,
				},
				{
					Name:   "Genre",
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
					Value:  anime["source"].(string),
					Inline: true,
				},
				{
					Name:   "Score",
					Value:  score,
					Inline: true,
				},
				{
					Name:   "Ranked",
					Value:  rank,
					Inline: true,
				},
				{
					Name:   "Popularity",
					Value:  strconv.Itoa(int(anime["popularity"].(float64))),
					Inline: true,
				},
				{
					Name:   "Members",
					Value:  strconv.Itoa(int(anime["members"].(float64))),
					Inline: true,
				},
			},
			Timestamp: time.Now().Format(time.RFC3339),
			Title:     title,
			URL:       anime["url"].(string),
			Image: &discordgo.MessageEmbedImage{
				URL: anime["image_url"].(string),
			},
		}
		_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	case <-time.After(config.Timeout * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, config.TimeoutMsg)
	}
	_ = ctx.Ses.ChannelMessageDelete(ctx.Msg.ChannelID, embedMsg.ID)
}
