package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"runtime"
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

	prefixMsg := "DM Prefix"
	prefix := config.Prefix
	if ctx.Msg.GuildID != "" {
		prefixMsg = "Server Prefix"
		prefix = framework.PDB.GetPrefix(ctx.Msg.GuildID)
	}

	// Memory Stats
	var memRuntime runtime.MemStats
	runtime.ReadMemStats(&memRuntime)

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Links",
				Value: fmt.Sprintf(
					"[Invite Bot](https://discordapp.com/oauth2/authorize?client_id=%s&scope=bot)\n"+
						"[Help Page](https://darenliang.github.io/MikuBot-Docs)\n"+
						"[Support Server](https://discord.gg/Tpa3cJB)\n"+
						"[Github Repo](https://github.com/darenliang/MikuBotGo)", client.ID),
				Inline: true,
			},
			{
				Name:   prefixMsg,
				Value:  prefix,
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
				Name:   "Servers",
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
			{
				Name:   "Heap In Use",
				Value:  fmt.Sprintf("%v MiB", memRuntime.HeapInuse/1024/1024),
				Inline: true,
			},
			{
				Name:   "Garbage Collected",
				Value:  fmt.Sprintf("%v GiB", memRuntime.HeapReleased/1024/1024/1024),
				Inline: true,
			},
			{
				Name:   "Go Version",
				Value:  runtime.Version(),
				Inline: true,
			},
		},
		Title: config.BotInfo,
	}
	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
