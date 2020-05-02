package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
	"path"
)

type NekoLifeResponse struct {
	URL string `json:"url"`
}

func getNekoLifeImage(tag string) (string, error) {
	nekoLifeResponse := NekoLifeResponse{}

	err := framework.UrlToStruct(fmt.Sprintf("https://nekos.life/api/v2/img/%s",
		tag), &nekoLifeResponse)

	if err != nil {
		log.Print("nekolife: response failed")
		return "", err
	}

	return nekoLifeResponse.URL, nil
}

// HeadPat command
func HeadPat(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("pat")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("headpat%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Here's your headpat.",
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

// Kiss command
func Kiss(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("kiss")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("kiss%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "OwO...",
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

// Tickle command
func Tickle(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("tickle")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("tickle%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Here's your tickle!",
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

// Feed command
func Feed(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("feed")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("feed%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Food. Yummy...",
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

// Slap command
func Slap(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("slap")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("slap%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Slap!",
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

// Cuddle command
func Cuddle(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("cuddle")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("cuddle%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "UwU...",
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

// Hug command
func Hug(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("hug")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("hug%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Here's your hug.",
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

// Smug command
func Smug(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("smug")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("smug%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Hehe...",
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

// Baka command
func Baka(ctx *exrouter.Context) {
	imageUrl, err := getNekoLifeImage("baka")

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
		return
	}

	resp, err := framework.HttpClient.Get(imageUrl)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot get image.")
		log.Printf("nekolife: failed to get image %s", imageUrl)
		return
	}

	fileName := fmt.Sprintf("baka%s", path.Ext(imageUrl))

	ms := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title: "Baka!",
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
