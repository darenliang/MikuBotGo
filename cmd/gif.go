package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/h2non/filetype"
	"github.com/mvdan/xurls"
	"io/ioutil"
	"net/http"
	"strings"
)

// Gif command
func Gif(ctx *exrouter.Context) {
	// Direct messages
	if len(ctx.Msg.GuildID) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "The Gif command cannot be used in DMs.")
		return
	}

	content := strings.TrimSpace(ctx.Args.After(1))

	if len(content) == 0 && len(ctx.Msg.Attachments) == 0 {
		title, link := framework.GBD.GetGif(ctx.Msg.GuildID)
		if title == "" && link == "" {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "There are no gifs currently stored in this guild.")
			return
		}
		usr, err := ctx.Ses.User(title)
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
			return
		}
		resp, err := http.Get(link)
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get gif from database.")
			return
		}
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf("Here's a gif from %s#%s",
			usr.Username, usr.Discriminator))
		_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, usr.ID+".gif", resp.Body)
		_ = resp.Body.Close()
		return
	}

	rxStrict := xurls.Strict
	urls := rxStrict.FindAllString(content, -1)

	gifUrls := make([]string, 0)
	for _, v := range urls {
		if strings.HasSuffix(v, ".gif") {
			gifUrls = append(gifUrls, v)
		}
	}
	for _, v := range ctx.Msg.Attachments {
		if strings.HasSuffix(v.URL, ".gif") {
			gifUrls = append(gifUrls, v.URL)
		}
	}

	if len(gifUrls) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You did not attach any gifs in the message.")
		return
	}

	count := 0

	// Filter fakes
	for _, v := range gifUrls {
		resp, err := http.Get(v)

		if err != nil {
			continue
		}

		data, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			continue
		}

		kind, _ := filetype.Match(data)
		if kind.Extension == "gif" && len(data) < config.MaxImgurByteSize {
			err := framework.GBD.UploadGif(ctx.Msg.GuildID, ctx.Msg.Author.ID, v)
			if err == nil {
				count++
			}
		}
		_ = resp.Body.Close()
	}

	msg := fmt.Sprintf("**%d** gif(s) added by %s#%s.",
		count, ctx.Msg.Author.Username, ctx.Msg.Author.Discriminator)

	if count != len(gifUrls) {
		msg += fmt.Sprintf(" **%d** gif(s) failed to upload.",
			len(gifUrls)-count)
	}

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, msg)
}
