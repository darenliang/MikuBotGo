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
	"os"
	"strconv"
	"strings"
)

var TraceMoeBase string

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

func init() {
	TraceMoeKey := os.Getenv("TRACEMOE")
	if TraceMoeKey != "" {
		TraceMoeBase = fmt.Sprintf("https://trace.moe/api/search?token=%s", TraceMoeKey)
	} else {
		TraceMoeBase = "https://trace.moe/api/search"
	}
}

func getPostJson(data string, target interface{}) error {
	form := url.Values{
		"image": {data},
	}

	r, err := framework.HttpClient.PostForm(TraceMoeBase, form)

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

// Identify command
func Identify(ctx *exrouter.Context) {
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
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%sidentify <anime screenshot (attachment or url)>`", prefix))
		return
	}

	if URL == "" {
		URL = ctx.Msg.Attachments[0].URL
	}

	ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer)

	defer ctx.Ses.MessageReactionRemove(ctx.Msg.ChannelID, ctx.Msg.ID, config.Timer, ctx.Ses.State.User.ID)

	image, err := framework.LoadImage(URL)

	if err != nil {
		ctx.Reply(":cry: There's an issue with your image.")
		return
	}

	image = imaging.Resize(image, 0, 480, imaging.Lanczos)

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, image, &jpeg.Options{Quality: 35})
	photoBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	trace := TraceData{}
	err = getPostJson(photoBase64, &trace)
	if err != nil {
		ctx.Reply(":cry: An error has occurred.")
		return
	}

	if len(trace.Docs) == 0 {
		ctx.Reply(":cry: We can't find any sauce from the provided image.")
		return
	}

	if trace.Docs[0].Similarity < 0.87 {
		ctx.Reply(":cry: We can't find the sauce for you.")
		return
	}

	response := framework.TraceSearchResult{}
	err = anilist.Query(framework.AnilistAnimeIDQuery(trace.Docs[0].AnilistID), &response)

	if err != nil {
		ctx.Reply(":cry: We can't get anime info of the sauce.")
		return
	}

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
				Name:   "Similarity",
				Value:  fmt.Sprintf("%.2f%%", trace.Docs[0].Similarity*100),
				Inline: true,
			},
			{
				Name:   "MAL Link",
				Value:  fmt.Sprintf("https://myanimelist.net/anime/%d", trace.Docs[0].MalID),
				Inline: false,
			},
		},
		Title: fmt.Sprintf("Here's the anime %s#%s", ctx.Msg.Author.Username, ctx.Msg.Author.Discriminator),
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

	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
