package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"time"
)

// Leaderboard command
func Leaderboard(ctx *exrouter.Context) {
	highScores := framework.GetHighscores()

	leaderboard := ""

	for idx, val := range highScores {
		if idx == 10 {
			break
		}
		leaderboard += fmt.Sprintf("`%2d.` ", idx+1)
		if idx == 0 {
			leaderboard += ":first_place:"
		} else if idx == 1 {
			leaderboard += ":second_place:"
		} else if idx == 2 {
			leaderboard += ":third_place:"
		} else {
			leaderboard += "   "
		}
		leaderboard += fmt.Sprintf(" `%3d` ", val.MusicScore)
		user, _ := ctx.Ses.User(val.UserId)
		leaderboard += fmt.Sprintf("%s#%s\n", user.Username, user.Discriminator)
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       config.EmbedColor,
		Description: leaderboard,
		Timestamp:   time.Now().Format(time.RFC3339),
		Title:       "Music Quiz Leaderboard",
	}

	_, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
}
