package main

import (
	"fmt"
	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/cmd"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
	"strings"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

// floppyEmoji
const floppyEmoji = "\xf0\x9f\x92\xbe"

func init() {
	waitReady := make(chan bool)
	var joinGuilds = make(map[string]bool)

	// Utility Group
	Router.Group(func(r *exrouter.Route) {
		// Info
		Router.OnMatch("info", dgrouter.NewRegexMatcher("^(?i)info$"), cmd.Info).Desc(
			"info: Get bot info\n\n" +
				"This command takes no arguments").Cat("Utility")

		// Prefix
		Router.OnMatch("prefix", dgrouter.NewRegexMatcher("^(?i)(prefix|p)$"), cmd.Prefix).Desc(
			"prefix: Set custom prefix\n\n" +
				"Alias: p\n\n" +
				"Usage: prefix <new prefix>\n\n" +
				"Note that you must have admin or owner privileges")

		// Ping
		Router.OnMatch("ping", dgrouter.NewRegexMatcher("^(?i)ping$"), cmd.Ping).Desc(
			"ping: Respond with pong\n\n" +
				"This command takes no arguments")

		// PFP command
		Router.OnMatch("pfp", dgrouter.NewRegexMatcher("^(?i)(pfp|avatar)$"), cmd.Pfp).Desc(
			"pfp: Get profile picture\n\n" +
				"Alias: avatar\n\n" +
				"Usage: pfp [@|username|username#tag|ID]")
	})

	// Anime Group
	Router.Group(func(r *exrouter.Route) {
		// Anime
		Router.OnMatch("anime", dgrouter.NewRegexMatcher("^(?i)(anime|a)$"), cmd.Anime).Desc(
			"anime: Get anime info\n\n" +
				"Alias: a\n\n" +
				"Usage: anime <anime name>").Cat("Anime")

		// Quiz
		Router.OnMatch("musicquiz", dgrouter.NewRegexMatcher("^(?i)(musicquiz|mq)$"), cmd.MusicQuiz).Desc(
			"musicquiz: Get an anime music quiz\n\n" +
				"Alias: mq\n\n" +
				"Usage:\n" +
				fmt.Sprintf("\t%-24v# Start an anime music quiz\n", "musicquiz") +
				fmt.Sprintf("\t%-24v# Guess an anime\n", "musicquiz <answer>") +
				fmt.Sprintf("\t%-24v# Give up current anime music quiz", "musicquiz giveup"))

		// Leaderboard
		Router.OnMatch("leaderboard", dgrouter.NewRegexMatcher("^(?i)(leaderboard|lb)$"), cmd.Leaderboard).Desc(
			"leaderboard: Get anime music leaderboard\n\n" +
				"Alias: lb\n\n" +
				"This command takes no arguments")

		// Trivia
		Router.OnMatch("trivia", dgrouter.NewRegexMatcher("^(?i)(trivia|t)$"), cmd.Trivia).Desc(
			"trivia: Get an anime trivia question\n\n" +
				"Alias: t\n\n" +
				"This command takes no arguments")

		// Waifu
		Router.OnMatch("waifu", dgrouter.NewRegexMatcher("^(?i)waifu$"), cmd.Waifu).Desc(
			"waifu: Get a never before seen waifu\n\n" +
				"Cross your fingers :)\n\n" +
				"This command takes no arguments")

		// Sauce
		Router.OnMatch("sauce", dgrouter.NewRegexMatcher("^(?i)sauce$"), cmd.Sauce).Desc(
			"sauce: Get sauce based on scene\n\n" +
				"Usage:\n" +
				"\tsauce <image url>\n" +
				"\tsauce <image attachment>\n")
	})

	// Fun Group
	Router.Group(func(r *exrouter.Route) {
		// Gif
		Router.OnMatch("gif", dgrouter.NewRegexMatcher("^(?i)gif$"), cmd.Gif).Desc(
			"gif: Add or get previously stored gifs\n\n" +
				"You can add gifs by reacting using the :floppy_disk: emoji\n\n" +
				"Usage:\n" +
				fmt.Sprintf("\t%-30v# Get a random gif\n", "gif") +
				fmt.Sprintf("\t%-30v# Store gifs\n", "gif <urls or attachments>")).Cat("Fun")
	})

	// Music Group
	Router.Group(func(r *exrouter.Route) {
		// Add
		Router.OnMatch("add", dgrouter.NewRegexMatcher("^(?i)add$"), cmd.AddMusic).Desc(
			"add: add music to queue\n\n" +
				"Usage:\n" +
				"\tplay <urls>\n").Cat("Music")

		// Clear
		Router.OnMatch("clear", dgrouter.NewRegexMatcher("^(?i)clear$"), cmd.ClearCommand).Desc(
			"clear: clear music in queue\n\n" +
				"This command takes no arguments")

		// Current
		Router.OnMatch("current", dgrouter.NewRegexMatcher("^(?i)current$"), cmd.CurrentCommand).Desc(
			"current: get current music\n\n" +
				"This command takes no arguments")

		// Join
		Router.OnMatch("join", dgrouter.NewRegexMatcher("^(?i)join$"), cmd.JoinCommand).Desc(
			"join: make bot join your voice channel\n\n" +
				"This command takes no arguments")

		// Leave
		Router.OnMatch("leave", dgrouter.NewRegexMatcher("^(?i)leave$"), cmd.LeaveCommand).Desc(
			"leave: make bot leave the voice channel\n\n" +
				"This command takes no arguments")

		// Pause
		Router.OnMatch("pause", dgrouter.NewRegexMatcher("^(?i)pause$"), cmd.PauseCommand).Desc(
			"pause: pause music\n\n" +
				"This command takes no arguments")

		// Play
		Router.OnMatch("play", dgrouter.NewRegexMatcher("^(?i)play$"), cmd.PlayCommand).Desc(
			"play: play music in queue\n\n" +
				"This command takes no arguments")

		// Queue
		Router.OnMatch("queue", dgrouter.NewRegexMatcher("^(?i)queue$"), cmd.QueueCommand).Desc(
			"queue: get music queue\n\n" +
				"This command takes no arguments")

		// Shuffle
		Router.OnMatch("shuffle", dgrouter.NewRegexMatcher("^(?i)shuffle$"), cmd.ShuffleCommand).Desc(
			"shuffle: shuffle music queue\n\n" +
				"This command takes no arguments")

		// Skip
		Router.OnMatch("skip", dgrouter.NewRegexMatcher("^(?i)skip$"), cmd.SkipCommand).Desc(
			"skip: skip current music\n\n" +
				"This command takes no arguments")

		// Stop
		Router.OnMatch("stop", dgrouter.NewRegexMatcher("^(?i)stop$"), cmd.StopCommand).Desc(
			"stop: stop music\n\n" +
				"This command takes no arguments")
	})

	// Help
	Router.Default = Router.OnMatch("help", dgrouter.NewRegexMatcher("^(?i)(help|h)$"), func(ctx *exrouter.Context) {
		command := strings.TrimSpace(ctx.Args.After(1))
		if command == "" {
			var text = fmt.Sprintf("help: Type %shelp <command> for more info on a command.\n",
				framework.PDB.GetPrefix(ctx.Msg.GuildID))
			pastCategory := ""
			for _, v := range Router.Routes {
				if v.Category != pastCategory && len(v.Category) != 0 {
					text += "\n" + v.Category + ":\n"
					pastCategory = v.Category
				}
				length := 16 - len(v.Name)
				text += "\t" + v.Name + strings.Repeat(" ", length) + "# " +
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
		"Usage: help <command>").Cat("Help")

	// Query database on ready
	Session.AddHandler(func(_ *discordgo.Session, ready *discordgo.Ready) {
		// Query databases to temp
		framework.PDB.SetGuilds()
		framework.MQDB.SetScores()
		framework.GBD.SetAlbums()

		// Music sessions and youtube
		config.MusicSessions = music.NewSessionManager()

		// Load cache and check for new guilds
		cache := framework.PDB.GetGuilds()
		for _, guild := range ready.Guilds {
			joinGuilds[guild.ID] = true
			if cache[guild.ID] == "" {
				framework.PDB.CreateGuild(guild.ID, config.Prefix)
			}
		}
		waitReady <- true
	})

	// Add guild on guild add
	Session.AddHandler(func(_ *discordgo.Session, create *discordgo.GuildCreate) {
		<-waitReady
		if !joinGuilds[create.ID] {
			framework.PDB.CreateGuild(create.ID, config.Prefix)
		}
	})

	// Remove guild on guild remote
	Session.AddHandler(func(_ *discordgo.Session, delete *discordgo.GuildDelete) {
		if !joinGuilds[delete.ID] {
			framework.PDB.RemoveGuild(delete.ID)
		}
	})

	// Handle incoming messages
	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := config.Prefix
		if len(m.GuildID) != 0 {
			prefix = framework.PDB.GetPrefix(m.GuildID)
		}
		_ = Router.FindAndExecute(Session, prefix, Session.State.User.ID, m.Message)
	})

	// Handle reaction add
	Session.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		// If DM
		if len(reaction.GuildID) == 0 {
			return
		}
		// If the emoji is not floppy disk
		if reaction.Emoji.Name != floppyEmoji {
			return
		}

		// Get message
		message, err := Session.ChannelMessage(reaction.ChannelID, reaction.MessageID)

		if err != nil {
			return
		}

		// Patch for ChannelMessage not returning guildID
		message.GuildID = reaction.GuildID

		// Get user
		user, err := Session.User(reaction.UserID)
		if err != nil {
			return
		}

		// Iterate emojis
		for _, emoji := range message.Reactions {
			if emoji.Emoji.Name == floppyEmoji && !emoji.Me {
				_ = Session.MessageReactionAdd(message.ChannelID, message.ID, floppyEmoji)
				count, total, dupCount, nsfwCount := cmd.UploadGifs(message.Content, message)
				if total == 0 {
					return
				}
				msg := cmd.GenerateGifUploadMessage(user, count, total, dupCount, nsfwCount)
				_, _ = Session.ChannelMessageSend(message.ChannelID, msg)
				return
			}
		}
	})
}
