package main

import (
	"fmt"
	"strings"

	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/cmd"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

func init() {

	// PFP command
	Router.On("pfp", cmd.Pfp).Desc(
		"pfp: Get profile pictures\n\n" +
			"Alias: avatar\n\n" +
			"Usage: [pfp|avatar] [@|username|username#tag|ID]",
	).Alias("avatar")

	// Ping
	Router.On("ping", cmd.Ping).Desc(
		"ping: Respond with pong\n\n" +
			"This command takes no arguments")

	// Info
	Router.On("info", cmd.Info).Desc(
		"info: Get bot info\n\n" +
			"This command takes no arguments")

	// Anime
	Router.On("anime", cmd.Anime).Desc(
		"anime: Get anime info\n\n" +
			"Alias: a\n\n" +
			"Usage: anime <anime name>").Alias("a")

	// Quiz
	Router.On("musicquiz", cmd.MusicQuiz).Desc(
		"musicquiz: Get an anime music quiz\n\n" +
			"Alias: mq\n\n" +
			"Usage:\n" +
			fmt.Sprintf("\t%-24v# Start an anime music quiz\n", "musicquiz") +
			fmt.Sprintf("\t%-24v# Guess an anime\n", "musicquiz <answer>") +
			fmt.Sprintf("\t%-24v# Give up current anime music quiz", "musicquiz giveup")).Alias("mq")

	// Leaderboard
	Router.On("leaderboard", cmd.Leaderboard).Desc(
		"leaderboard: Get anime music leaderboard\n\n" +
			"Alias: lb\n\n" +
			"This command takes no arguments").Alias("lb")

	// Help
	Router.Default = Router.On("help", func(ctx *exrouter.Context) {
		command := strings.TrimSpace(ctx.Args.After(1))

		if command == "" {
			var text = fmt.Sprintf("help: Type %shelp <command> for more info on a command.\n\n", config.Prefix)
			for _, v := range Router.Routes {
				length := 16 - len(v.Name)
				text += v.Name + strings.Repeat(" ", length) + "# " +
					framework.GeneratePreviewDesc(v.Description) + "\n"
			}
			_, _ = ctx.Reply("```\n" + text + "```\n")
			return
		}

		for _, v := range Router.Routes {
			if command == v.Name {
				_, _ = ctx.Reply("```\n" + v.Description + "```\n")
				return
			}
		}

		_, _ = ctx.Reply("Command not found.")
	}).Desc("help: Prints this help menu\n\n" +
		"Alias: h\n\n" +
		"Usage: help <command>").Alias("h")

	// Handle incoming messages
	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		_ = Router.FindAndExecute(Session, config.Prefix, Session.State.User.ID, m.Message)
	})
}
