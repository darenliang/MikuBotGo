package main

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/cmd"
	"github.com/darenliang/MikuBotGo/config"
	"strings"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

func init() {
	// Ping
	Router.On("ping", cmd.Ping).Desc("responds with pong")

	// Info
	Router.On("info", cmd.Info).Desc("get bot info")

	// Anime
	Router.On("anime", cmd.Anime).Desc("get anime info")

	// Quiz
	Router.On("musicquiz", cmd.MusicQuiz).Desc("get an anime music quiz")

	// Leaderboard
	Router.On("leaderboard", cmd.Leaderboard).Desc("get anime music leaderboard")

	// Help
	Router.Default = Router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range Router.Routes {
			length := 16 - len(v.Name)
			text += v.Name + strings.Repeat(" ", length) + "# " + v.Description + "\n"
		}
		_, _ = ctx.Reply("```\n" + text + "```\n")
	}).Desc("prints this help menu")

	// Handle incoming messages
	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		_ = Router.FindAndExecute(Session, config.Prefix, Session.State.User.ID, m.Message)
	})
}
