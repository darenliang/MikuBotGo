package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"sort"
)

// Leaderboard command
func Leaderboard(ctx *exrouter.Context) {
	scores := framework.MQDB.GetScores()

	scoresSlice := make([]framework.MusicQuizEntry, 0, len(scores))
	for k, v := range scores {
		scoresSlice = append(scoresSlice, framework.MusicQuizEntry{
			UserId:        k,
			MusicScore:    v.MusicScore,
			TotalAttempts: v.TotalAttempts,
		})
	}

	sort.Slice(scoresSlice, func(i, j int) bool {
		return scoresSlice[i].MusicScore > scoresSlice[j].MusicScore
	})

	leaderboard := "```\nRank | Score | User\n"

	for idx, val := range scoresSlice {
		if idx == 10 {
			break
		}
		leaderboard += fmt.Sprintf("%4d |", idx+1)
		leaderboard += fmt.Sprintf(" %5d |", val.MusicScore*100)
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
