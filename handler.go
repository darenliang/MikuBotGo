package main

import (
	"github.com/Necroforger/dgrouter"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/cmd"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"github.com/darenliang/MikuBotGo/music"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

// floppyEmoji
const floppyEmoji = "\xf0\x9f\x92\xbe"

func init() {
	waitReady := make(chan bool)
	var joinGuilds = make(map[string]bool)

	Router.OnMatch("info", dgrouter.NewRegexMatcher("^(?i)info$"), cmd.Info)

	Router.OnMatch("prefix", dgrouter.NewRegexMatcher("^(?i)(prefix|p)$"), cmd.Prefix)

	Router.OnMatch("pfp", dgrouter.NewRegexMatcher("^(?i)(pfp|avatar)$"), cmd.Pfp)

	Router.OnMatch("cloc", dgrouter.NewRegexMatcher("^(?i)cloc$"), cmd.Cloc)

	Router.OnMatch("anime", dgrouter.NewRegexMatcher("^(?i)(anime|a)$"), cmd.Anime)

	Router.OnMatch("manga", dgrouter.NewRegexMatcher("^(?i)(manga|m)$"), cmd.Manga)

	Router.OnMatch("musicquiz", dgrouter.NewRegexMatcher("^(?i)(musicquiz|mq)$"), cmd.MusicQuiz)

	Router.OnMatch("leaderboard", dgrouter.NewRegexMatcher("^(?i)(leaderboard|lb)$"), cmd.Leaderboard)

	Router.OnMatch("trivia", dgrouter.NewRegexMatcher("^(?i)(trivia|t)$"), cmd.Trivia)

	Router.OnMatch("waifu", dgrouter.NewRegexMatcher("^(?i)waifu$"), cmd.Waifu)

	Router.OnMatch("sauce", dgrouter.NewRegexMatcher("^(?i)sauce$"), cmd.Sauce)

	Router.OnMatch("gif", dgrouter.NewRegexMatcher("^(?i)gif$"), cmd.Gif)

	Router.OnMatch("owofy", dgrouter.NewRegexMatcher("^(?i)owofy$"), cmd.Owofy)

	Router.OnMatch("8ball", dgrouter.NewRegexMatcher("^(?i)8ball|8b|eightball$"), cmd.EightBall)

	Router.OnMatch("kaomoji", dgrouter.NewRegexMatcher("^(?i)kaomoji$"), cmd.Kaomoji)

	Router.OnMatch("add", dgrouter.NewRegexMatcher("^(?i)add$"), cmd.AddMusic)

	Router.OnMatch("clear", dgrouter.NewRegexMatcher("^(?i)clear$"), cmd.ClearCommand)

	Router.OnMatch("current", dgrouter.NewRegexMatcher("^(?i)current$"), cmd.CurrentCommand)

	Router.OnMatch("join", dgrouter.NewRegexMatcher("^(?i)join$"), cmd.JoinCommand)

	Router.OnMatch("leave", dgrouter.NewRegexMatcher("^(?i)leave$"), cmd.LeaveCommand)

	Router.OnMatch("pause", dgrouter.NewRegexMatcher("^(?i)pause$"), cmd.PauseCommand)

	Router.OnMatch("play", dgrouter.NewRegexMatcher("^(?i)play$"), cmd.PlayCommand)

	Router.OnMatch("queue", dgrouter.NewRegexMatcher("^(?i)queue$"), cmd.QueueCommand)

	Router.OnMatch("shuffle", dgrouter.NewRegexMatcher("^(?i)shuffle$"), cmd.ShuffleCommand)

	Router.OnMatch("skip", dgrouter.NewRegexMatcher("^(?i)skip$"), cmd.SkipCommand)

	Router.OnMatch("stop", dgrouter.NewRegexMatcher("^(?i)stop$"), cmd.StopCommand)

	Router.Default = Router.OnMatch("help", dgrouter.NewRegexMatcher("^(?i)(help|h)$"), func(ctx *exrouter.Context) {
		_, _ = ctx.Reply("Please visit __https://darenliang.github.io/MikuBot-Docs__ for help on all the commands.")
	}).Cat("Help")

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
