package main

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/cmd"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"strings"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

func init() {

	// PFP command
	Router.On("pfp", cmd.Pfp).Desc(
		"pfp: Get profile picture\n\n" +
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

	// Prefix
	Router.On("prefix", cmd.Prefix).Desc(
		"prefix: Set custom prefix\n\n" +
			"Alias: p\n\n" +
			"Usage: prefix <new prefix>\n\n" +
			"Note that you must have admin or owner privileges").Alias("p")

	// Quiz
	Router.On("musicquiz", cmd.MusicQuiz).Desc(
		"musicquiz: Get an anime music quiz\n\n" +
			"Alias: mq\n\n" +
			"Usage:\n" +
			fmt.Sprintf("\t%-24v# Start an anime music quiz\n", "musicquiz") +
			fmt.Sprintf("\t%-24v# Guess an anime\n", "musicquiz <answer>") +
			fmt.Sprintf("\t%-24v# Give up current anime music quiz", "musicquiz giveup")).Alias("mq")

	// Trivia
	Router.On("trivia", cmd.Trivia).Desc(
		"trivia: Get an anime trivia question\n\n" +
			"Alias: t\n\n" +
			"This command takes no arguments").Alias("t")

	// Leaderboard
	Router.On("leaderboard", cmd.Leaderboard).Desc(
		"leaderboard: Get anime music leaderboard\n\n" +
			"Alias: lb\n\n" +
			"This command takes no arguments").Alias("lb")

	// Sauce
	Router.On("sauce", cmd.Sauce).Desc(
		"sauce: Get sauce based on scene\n\n" +
			"Usage:\n" +
			"\tsauce <image url>\n" +
			"\tsauce <image attachment>\n")

	// Help
	Router.Default = Router.On("help", func(ctx *exrouter.Context) {
		command := strings.TrimSpace(ctx.Args.After(1))
		if command == "" {
			var text = fmt.Sprintf("help: Type %shelp <command> for more info on a command.\n\n",
				framework.PDB.GetPrefix(ctx.Msg.GuildID))
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

	// Query database on ready
	Session.AddHandler(func(_ *discordgo.Session, ready *discordgo.Ready) {
		// Query database to temp
		framework.PDB.SetGuilds()

		// Load cache and check for new guilds
		cache := framework.PDB.GetGuilds()
		for _, guild := range ready.Guilds {
			if cache[guild.ID] == "" {
				framework.PDB.CreateGuild(guild.ID, config.Prefix)
			}
		}
	})

	// Add guild on guild add
	Session.AddHandler(func(_ *discordgo.Session, create *discordgo.GuildCreate) {
		framework.PDB.CreateGuild(create.ID, config.Prefix)
	})

	// Remove guild on guild remote
	Session.AddHandler(func(_ *discordgo.Session, delete *discordgo.GuildDelete) {
		framework.PDB.RemoveGuild(delete.ID)
	})

	// Handle incoming messages
	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := config.Prefix
		if len(m.GuildID) != 0 {
			prefix = framework.PDB.GetPrefix(m.GuildID)
		}
		_ = Router.FindAndExecute(Session, prefix, Session.State.User.ID, m.Message)
	})
}
