package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Waifu command
func Waifu(ctx *exrouter.Context) {
	imageID := rand.Int() % 100001

	resp, err := framework.HttpClient.Get(fmt.Sprintf(
		"https://www.thiswaifudoesnotexist.net/example-%d.jpg", imageID))

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot generate an image.")
		log.Printf("waifu: generate image error with id %d", imageID)
		return
	}

	fileName := fmt.Sprintf("waifu.jpg")

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Here's your waifu.",
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

	_, _ = ctx.Ses.ChannelMessageSendComplex(ctx.Msg.ChannelID, ms)
}
