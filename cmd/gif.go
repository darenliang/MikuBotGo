package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/h2non/filetype"
	"github.com/mvdan/xurls"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ClarifaiPredict struct {
	Status struct {
		Code int `json:"code"`
	} `json:"status"`
	Outputs []struct {
		Data struct {
			Frames []struct {
				Data struct {
					Concepts []struct {
						Name  string  `json:"name"`
						Value float64 `json:"value"`
					} `json:"concepts"`
				} `json:"data"`
			} `json:"frames"`
		} `json:"data"`
	} `json:"outputs"`
}

func moderateGif(url string) (bool, error) {
	jsonStr := fmt.Sprintf(`{
   "inputs":[
      {
         "data":{
            "video":{
               "url":"%s"
            }
         }
      }
   ],
   "model":{
      "output_info":{
         "output_config":{
            "sample_ms":100
         }
      }
   }
}`, url)
	req, _ := http.NewRequest("POST", config.ClarifaiNSFWEndpoint, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+config.ClarifaiToken)
	resp, _ := framework.HttpClient.Do(req)

	clarifaiPredict := ClarifaiPredict{}
	_ = json.NewDecoder(resp.Body).Decode(&clarifaiPredict)
	if clarifaiPredict.Status.Code != 10000 || len(clarifaiPredict.Outputs) == 0 {
		return false, errors.New("invalid status code")
	}
	for _, frame := range clarifaiPredict.Outputs[0].Data.Frames {
		if frame.Data.Concepts[0].Name == "sfw" {
			continue
		}
		if frame.Data.Concepts[0].Value > 0.9 {
			return false, nil
		}
	}
	return true, nil
}

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
	nsfwCount := 0
	dupCount := 0

	// Filter fakes
	for _, v := range gifUrls {
		resp, err := http.Get(v)

		if err != nil {
			continue
		}

		data, err := ioutil.ReadAll(io.LimitReader(resp.Body, config.MaxImgurByteSize+1))

		// URL read is good
		if err == nil {
			// Size too big
			if len(data) <= config.MaxImgurByteSize {
				// Validate filetype
				kind, _ := filetype.Match(data)
				if kind.Extension == "gif" {
					// Moderate file
					ok, err := moderateGif(v)
					if err == nil {
						if ok {
							// Check file hash for dups
							hash := fmt.Sprintf("%x", sha256.Sum256(data))
							if framework.GBD.CheckDup(ctx.Msg.GuildID, hash) {
								dupCount++
							} else {
								// Update file
								err := framework.GBD.UploadGif(ctx.Msg.GuildID, ctx.Msg.Author.ID, v, hash)
								if err == nil {
									count++
								}
							}
						} else {
							// If failed moderation
							nsfwCount++
						}
					}
				}
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

	if dupCount != 0 {
		msg += fmt.Sprintf(" **%d** gif(s) was flagged as duplicate.",
			dupCount)
	}

	if nsfwCount != 0 {
		msg += fmt.Sprintf(" **%d** gif(s) was flagged as questionable.",
			nsfwCount)
	}

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, msg)
}
