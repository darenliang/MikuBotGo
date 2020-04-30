package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"net/url"
	"path"
)

type DanbooruResponse []struct {
	LargeFileURL string `json:"large_file_url"`
}

// Cat girl command
func CatGirl(ctx *exrouter.Context) {
	danbooru := DanbooruResponse{}

	err := framework.UrlToStruct(fmt.Sprintf("https://danbooru.donmai.us/posts.json?login=%s&api_key=%s&random=true&limit=1&tags=%s",
		config.DanbooruUsername, config.DanbooruToken, url.QueryEscape("cat_girl score:>=35 rating:safe")), &danbooru)

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		log.Print("catgirl: response failed")
		return
	}

	if len(danbooru) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		log.Print("catgirl: response empty")
		return
	}

	resp, err := framework.HttpClient.Get(danbooru[0].LargeFileURL)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		log.Print("catgirl: failed to get image")
		return
	}

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Here's your catgirl.")
	_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, fmt.Sprintf("catgirl%s", path.Ext(danbooru[0].LargeFileURL)), resp.Body)
}
