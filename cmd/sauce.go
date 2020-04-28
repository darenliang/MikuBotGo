package cmd

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/animenotifier/anilist"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"net/url"
	"strconv"
	"strings"
)

type TraceData struct {
	Docs []struct {
		AnilistID  int         `json:"anilist_id"`
		At         float64     `json:"at"`
		Filename   string      `json:"filename"`
		Episode    interface{} `json:"episode,string"`
		Tokenthumb string      `json:"tokenthumb"`
		Similarity float64     `json:"similarity"`
		MalID      int         `json:"mal_id"`
	} `json:"docs"`
}

func getPostJson(data string, target interface{}) error {
	form := url.Values{
		"image": {data},
	}

	r, err := framework.HttpClient.PostForm(config.TraceMoeBase, form)

	if r != nil {
		defer r.Body.Close()
	}

	if err != nil {
		return err
	}

	err = json.NewDecoder(r.Body).Decode(target)
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

	_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	image, err := framework.LoadImage(URL)

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Theres an issue with your image.")
		_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
	}

	image = imaging.Resize(image, 0, 480, imaging.Lanczos)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, image, &jpeg.Options{Quality: 35})
	photoBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	trace := TraceData{}
	err = getPostJson(photoBase64, &trace)
	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
		return
	}

	if len(trace.Docs) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "We can't find any sauce from the provided image.")
		_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
		return
	}

	if trace.Docs[0].Similarity < 0.87 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "We can't find the sauce for you.")
		_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
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
			URL: fmt.Sprintf("https://trace.moe/thumbnail.php?anilist_id=%d&file=%s&t=%g&token=%s",
				trace.Docs[0].AnilistID,
				url.QueryEscape(trace.Docs[0].Filename),
				trace.Docs[0].At,
				trace.Docs[0].Tokenthumb),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: anime.CoverImage.ExtraLarge,
		},
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)

	_ = ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)
}
