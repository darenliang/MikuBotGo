package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"strconv"
	"time"
)

// Info command
func Info(ctx *exrouter.Context) {
	client, _ := ctx.Ses.Application("@me")
	currTime := time.Since(config.StartTime)

	userCount := 0
	for _, guild := range ctx.Ses.State.Guilds {
		userCount += guild.MemberCount
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Links",
				Value: fmt.Sprintf(
					"[Invite Bot](https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot)\n"+
						"[Support Server](https://discord.gg/Tpa3cJB)\n"+
						"[Github Repo](https://github.com/darenliang/MikuBotGo)", client.ID),
				Inline: true,
			},
			{
				Name:   "Server Prefix",
				Value:  config.Prefix,
				Inline: true,
			},
			{
				Name:   "Created by",
				Value:  fmt.Sprintf("%s#%s", client.Owner.Username, client.Owner.Discriminator),
				Inline: false,
			},
			{
				Name:   "Latency",
				Value:  fmt.Sprintf("%dms", ctx.Ses.HeartbeatLatency().Milliseconds()),
				Inline: true,
			},
			{
				Name:   "Guilds",
				Value:  fmt.Sprintf("%d", len(ctx.Ses.State.Guilds)),
				Inline: true,
			},
			{
				Name:   "Users",
				Value:  fmt.Sprintf("%d", userCount),
				Inline: true,
			},
			{
				Name: "Uptime",
				Value: fmt.Sprintf("%dd, %dh, %dm, %ds",
					int(currTime.Hours())/24,
					int(currTime.Hours())%24,
					int(currTime.Minutes())%60,
					int(currTime.Seconds())%60),
				Inline: true,
			},
			{
				Name:   "Current Shard ID",
				Value:  strconv.Itoa(ctx.Ses.ShardID),
				Inline: true,
			},
			{
				Name:   "Shard Count",
				Value:  strconv.Itoa(ctx.Ses.ShardCount),
				Inline: true,
			},
		},
		Title: config.BotInfo,
	}
	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
