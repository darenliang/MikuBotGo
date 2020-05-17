package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"sort"
	"strings"
)

// Leaderboard command
func Leaderboard(ctx *exrouter.Context) {
	query := strings.TrimSpace(ctx.Args.After(1))
	scores := framework.MQDB.GetScores()

	memberSet := make(map[string]bool)
	var memberFilter func(id string) bool

	title := ""

	if query == "global" || query == "g" || ctx.Msg.GuildID == "" {
		title = "Global Music Quiz Leaderboard"

		memberFilter = func(id string) bool {
			return true
		}
	} else {
		guild, _ := ctx.Guild(ctx.Msg.GuildID)

		title = fmt.Sprintf("Music Quiz Leaderboard for %s", guild.Name)

		for _, member := range guild.Members {
			memberSet[member.User.ID] = true
		}

		memberFilter = func(id string) bool {
			return memberSet[id]
		}
	}

	scoresSlice := make([]framework.MusicQuizEntry, 0, len(scores))
	for k, v := range scores {
		if !memberFilter(k) {
			continue
		}
		scoresSlice = append(scoresSlice, framework.MusicQuizEntry{
			UserId:        k,
			MusicScore:    v.MusicScore,
			TotalAttempts: v.TotalAttempts,
		})
	}

	sort.Slice(scoresSlice, func(i, j int) bool {
		return scoresSlice[i].MusicScore > scoresSlice[j].MusicScore
	})

	leaderboard := "```\nRank |  Score | User\n"

	for idx, val := range scoresSlice {
		if idx == 10 || val.MusicScore == 0 {
			break
		}
		leaderboard += fmt.Sprintf("%4d |", idx+1)
		leaderboard += fmt.Sprintf(" %6d |", val.MusicScore*100)
		user, _ := ctx.Ses.User(val.UserId)
		leaderboard += fmt.Sprintf(" %s#%s\n", user.Username, user.Discriminator)
	}

	leaderboard += "\n```"

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: leaderboard,
		Title:       title,
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
