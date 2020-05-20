package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"log"
)

// Pfp command
func Pfp(ctx *exrouter.Context) {

	target := framework.Getuser(ctx)

	if target != nil {

		resp, err := framework.HttpClient.Get(target.AvatarURL("1024"))

		if resp != nil {
			defer resp.Body.Close()
		}

		if err != nil {
			ctx.Reply(":cry: Failed to get profile pic.")
			log.Print("pfp: failed to get image")
			return
		}

		fileName := fmt.Sprintf("profile.png")

		ms := &discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: fmt.Sprintf("Here's %s#%s's profile pic.", target.Username, target.Discriminator),
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
	} else {
		ctx.Reply(":information_source: The user is not found.")
	}

}
