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
	"os"
	"strings"
)

const (
	// 10 MB is the max imgur size
	maxImgurByteSize = 1000 * 1000 * 10
	// Clarifai API Endpoint
	clarifaiNSFWEndpoint = "https://api.clarifai.com/v2/models/e9576d86d2004ed1a38ba0cf39ecb4b1/outputs"
)

var clarifaiToken string

func init() {
	clarifaiToken = os.Getenv("CLARIFAI_TOKEN")
}

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
	req, _ := http.NewRequest("POST", clarifaiNSFWEndpoint, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Key "+clarifaiToken)
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
			log.Printf("gif: url fail: %s", v)
			continue
		}

		data, err := ioutil.ReadAll(io.LimitReader(resp.Body, maxImgurByteSize+1))

		// URL read is good
		if err == nil {
			// Size too big
			if len(data) <= maxImgurByteSize {
				// Validate filetype
				kind, _ := filetype.Match(data)
				if kind.Extension == "gif" {
					// Moderate file
					ok, err := moderateGif(v)

					if err == nil && !ok {
						log.Printf("gif: moderation fail: %s", v)
						nsfwCount++
					} else {
						// Check file hash for dups
						hash := fmt.Sprintf("%x", sha256.Sum256(data))
						if framework.GBD.CheckDup(message.GuildID, hash) {
							log.Print("gif: dup found")
							dupCount++
						} else {
							// Update file
							err := framework.GBD.UploadGif(message.GuildID, message.Author.ID, v, hash)
							if err == nil {
								count++
							} else {
								log.Printf("gif: upload fail: %s", v)
							}
						}
					}
				} else {
					log.Printf("gif: filekind not supported: %s", v)
				}
			} else {
				log.Printf("gif: too large: %s", v)
			}
		} else {
			log.Printf("gif: readall failed: %s", v)
		}
		resp.Body.Close()
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
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	// Direct messages
	if ctx.Msg.GuildID == "" {
		ctx.Reply(":warning: The gif command cannot be used in DMs.")
		return
	}

	content := strings.TrimSpace(ctx.Args.After(1))

	if len(content) == 0 && len(ctx.Msg.Attachments) == 0 {
		title, link := framework.GBD.GetGif(ctx.Msg.GuildID)
		if title == "" && link == "" {
			ctx.Reply(fmt.Sprintf(":cry: There are no gifs currently stored for this server. To store gifs use the command `%sgif <gif attachment(s) or url(s)>`", prefix))
			return
		}
		usr, err := ctx.Ses.User(title)
		if err != nil {
			ctx.Reply(":cry: Cannot get user data.")
			log.Printf("gif: user not found")
			return
		}
		resp, err := framework.HttpClient.Get(link)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			ctx.Reply(":cry: Cannot get gif from database.")
			log.Printf("gif: gif link errored out")
			return
		}

		fileName := usr.ID + ".gif"

		ms := &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: fmt.Sprintf("Here's a gif from %s#%s.",
					usr.Username, usr.Discriminator),
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

		return
	}

	count, total, dupCount, nsfwCount := UploadGifs(content, ctx.Msg)

	msg := GenerateGifUploadMessage(ctx.Msg.Author, count, total, dupCount, nsfwCount)

	ctx.Reply(msg)
}
