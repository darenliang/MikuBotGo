package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
)

// Leaderboard command
func Leaderboard(ctx *exrouter.Context) {
	highScores := framework.GetHighscores()

	leaderboard := "```\nRank | Score | Invokes | User\n"

	for idx, val := range highScores {
		if idx == 10 {
			break
		}
		leaderboard += fmt.Sprintf("%4d |", idx+1)
		leaderboard += fmt.Sprintf(" %5d |", val.MusicScore*100)
		leaderboard += fmt.Sprintf(" %7d |", val.TotalAttempts)
		user, _ := ctx.Ses.User(val.UserId)
		leaderboard += fmt.Sprintf(" %s#%s\n", user.Username, user.Discriminator)
	}

	leaderboard += "\n```"

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: leaderboard,
		Title:       "Music Quiz Leaderboard",
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
