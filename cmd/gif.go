package cmd

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/h2non/filetype"
	"github.com/mvdan/xurls"
	"io"
	"io/ioutil"
	"log"
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
	resp, err := framework.HttpClient.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		log.Printf("gif: clarifai api call failed")
		return false, nil
	}

	clarifaiPredict := ClarifaiPredict{}
	err = json.NewDecoder(resp.Body).Decode(&clarifaiPredict)

	if err != nil {
		log.Printf("gif: clarifai response decode failed")
		return false, err
	}

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

func UploadGifs(content string, message *discordgo.Message) (int, int, int, int) {
	count := 0
	dupCount := 0
	nsfwCount := 0

	rxStrict := xurls.Strict
	urls := rxStrict.FindAllString(content, -1)

	gifUrls := make([]string, 0)
	for _, v := range urls {
		if strings.HasSuffix(v, ".gif") {
			gifUrls = append(gifUrls, v)
		}
	}
	for _, v := range message.Attachments {
		if strings.HasSuffix(v.URL, ".gif") {
			gifUrls = append(gifUrls, v.URL)
		}
	}

	// Filter fakes
	for _, v := range gifUrls {
		resp, err := framework.HttpClient.Get(v)

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
							if framework.GBD.CheckDup(message.GuildID, hash) {
								dupCount++
							} else {
								// Update file
								err := framework.GBD.UploadGif(message.GuildID, message.Author.ID, v, hash)
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
		} else {
			log.Printf("gif: ioutil readall failed")
		}
		_ = resp.Body.Close()
	}

	return count, len(gifUrls), dupCount, nsfwCount
}

func GenerateGifUploadMessage(user *discordgo.User, count, total, dupCount, nsfwCount int) string {
	msg := fmt.Sprintf("**%d** gif(s) added by %s#%s.",
		count, user.Username, user.Discriminator)

	if count != total {
		msg += fmt.Sprintf(" **%d** gif(s) failed to upload.",
			total-count)
	}

	if dupCount != 0 {
		msg += fmt.Sprintf(" **%d** gif(s) was flagged as duplicate.",
			dupCount)
	}

	if nsfwCount != 0 {
		msg += fmt.Sprintf(" **%d** gif(s) was flagged as questionable.",
			nsfwCount)
	}

	return msg
}

// Gif command
func Gif(ctx *exrouter.Context) {
	// Direct messages
	if ctx.Msg.GuildID == "" {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "The Gif command cannot be used in DMs.")
		return
	}

	content := strings.TrimSpace(ctx.Args.After(1))

	if len(content) == 0 && len(ctx.Msg.Attachments) == 0 {
		title, link := framework.GBD.GetGif(ctx.Msg.GuildID)
		if title == "" && link == "" {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "There are no gifs currently stored for this server.")
			return
		}
		usr, err := ctx.Ses.User(title)
		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
			log.Printf("gif: user not found")
			return
		}
		resp, err := framework.HttpClient.Get(link)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get gif from database.")
			log.Printf("gif: gif link errored out")
			return
		}

		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf("Here's a gif from %s#%s",
			usr.Username, usr.Discriminator))
		_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, usr.ID+".gif", resp.Body)
		return
	}

	count, total, dupCount, nsfwCount := UploadGifs(content, ctx.Msg)

	msg := GenerateGifUploadMessage(ctx.Msg.Author, count, total, dupCount, nsfwCount)

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, msg)
}
