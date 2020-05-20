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

var StartTime time.Time

func init() {
	// Set start time
	StartTime = time.Now()
}

// Info command
func Info(ctx *exrouter.Context) {
	client, _ := ctx.Ses.Application("@me")
	currTime := time.Since(StartTime)

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

	// Last GC
	lastGC := time.Since(time.Unix(0, int64(memRuntime.LastGC)))

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: "(づ｡◕‿‿◕｡)づ Made with :heart: and DiscordGo.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name: "Links",
				Value: fmt.Sprintf(
					"[Invite Bot](https://discord.com/oauth2/authorize?client_id=%s&scope=bot)\n"+
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
				Name:   "Memory Alloc",
				Value:  fmt.Sprintf("%v MiB", memRuntime.HeapAlloc/1024/1024),
				Inline: true,
			},
			{
				Name:   "Memory Target",
				Value:  fmt.Sprintf("%v MiB", memRuntime.NextGC/1024/1024),
				Inline: true,
			},
			{
				Name: "Last GC",
				Value: fmt.Sprintf("%dm, %ds",
					int(lastGC.Minutes())%60,
					int(lastGC.Seconds())%60),
				Inline: true,
			},
			{
				Name:   "Available Cores",
				Value:  fmt.Sprintf("%v", runtime.NumCPU()),
				Inline: true,
			},
			{
				Name:   "Subroutines",
				Value:  fmt.Sprintf("%v", runtime.NumGoroutine()),
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
	ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
