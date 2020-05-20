package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

type DanbooruResponse []struct {
	LargeFileURL string `json:"large_file_url"`
}

var (
	DanbooruUsername string
	DanbooruToken    string
)

func init() {
	DanbooruUsername = os.Getenv("DANBOORU_USERNAME")
	DanbooruToken = os.Getenv("DANBOORU_TOKEN")
}

// Cat girl command
func CatGirl(ctx *exrouter.Context) {
	danbooru := DanbooruResponse{}

	err := framework.UrlToStruct(fmt.Sprintf("https://danbooru.donmai.us/posts.json?login=%s&api_key=%s&random=true&limit=10&tags=%s",
		DanbooruUsername, DanbooruToken, url.QueryEscape("cat_girl score:>=35 rating:safe")), &danbooru)

	if err != nil {
		ctx.Reply(":cry: Sorry, failed to get catgirl.")
		log.Print("catgirl: response failed")
		return
	}

	if len(danbooru) == 0 {
		ctx.Reply(":cry: Sorry, failed to get catgirl.")
		log.Print("catgirl: response empty")
		return
	}

	var fileUrl string

	for _, v := range danbooru {
		if strings.HasSuffix(v.LargeFileURL, ".jpg") || strings.HasSuffix(v.LargeFileURL, ".png") || strings.HasSuffix(v.LargeFileURL, ".gif") {
			fileUrl = v.LargeFileURL
			break
		}
	}

	if fileUrl == "" {
		ctx.Reply(":cry: Sorry, failed to get catgirl.")
		log.Print("catgirl: no suitable image")
		return
	}

	resp, err := framework.HttpClient.Get(danbooru[0].LargeFileURL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		ctx.Reply(":cry: Sorry, failed to get catgirl.")
		log.Printf("catgirl: failed to get image: %s", danbooru[0].LargeFileURL)
		return
	}

	fileName := fmt.Sprintf("catgirl%s", path.Ext(danbooru[0].LargeFileURL))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Here's your catgirl.",
			Color: config.EmbedColor,
			Image: &discordgo.MessageEmbedImage{
				URL: "attachment://" + fileName,
			},
		},
		Files: []*discordgo.File{
			{
				Name:   fileName,
				Reader: resp.Body,
			},
		},
	}

	ctx.Ses.ChannelMessageSendComplex(ctx.Msg.ChannelID, ms)
}
