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
	"github.com/foxbot/gavalink"
	"log"
	"sync"
)

// Router is registered as a global variable to allow easy access to the
// multiplexer throughout the bot.
var Router = exrouter.New()

// floppyEmoji
const floppyEmoji = "\xf0\x9f\x92\xbe"

// waitReady is a thread safe ready signal
type waitReady struct {
	ready bool
	mux   sync.Mutex
}

// setReady sets a thread safe ready signal
func (c *waitReady) setReady() {
	c.mux.Lock()
	c.ready = true
	c.mux.Unlock()
}

// setReady gets a thread safe ready signal
func (c *waitReady) getReady() bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.ready
}

func init() {
	// Initialize ready signal
	status := waitReady{
		ready: false,
	}

	// Create map of existing guilds
	var readyGuilds = make(map[string]bool)

	// Add command handlers
	// Utility
	Router.OnMatch("info", dgrouter.NewRegexMatcher("^(?i)info$"), cmd.Info)
	Router.OnMatch("prefix", dgrouter.NewRegexMatcher("^(?i)(prefix|p)$"), cmd.Prefix)
	Router.OnMatch("pfp", dgrouter.NewRegexMatcher("^(?i)(pfp|avatar)$"), cmd.Pfp)
	Router.OnMatch("cloc", dgrouter.NewRegexMatcher("^(?i)cloc$"), cmd.Cloc)

	// Anime
	Router.OnMatch("anime", dgrouter.NewRegexMatcher("^(?i)(anime|a)$"), cmd.Anime)
	Router.OnMatch("manga", dgrouter.NewRegexMatcher("^(?i)(manga|m)$"), cmd.Manga)
	Router.OnMatch("musicquiz", dgrouter.NewRegexMatcher("^(?i)(musicquiz|mq)$"), cmd.MusicQuiz)
	Router.OnMatch("leaderboard", dgrouter.NewRegexMatcher("^(?i)(leaderboard|lb)$"), cmd.Leaderboard)
	Router.OnMatch("trivia", dgrouter.NewRegexMatcher("^(?i)(trivia|t)$"), cmd.Trivia)
	Router.OnMatch("waifu", dgrouter.NewRegexMatcher("^(?i)waifu$"), cmd.Waifu)
	Router.OnMatch("sauce", dgrouter.NewRegexMatcher("^(?i)sauce$"), cmd.Sauce)
	Router.OnMatch("identify", dgrouter.NewRegexMatcher("^(?i)identify$"), cmd.Identify)

	// Fun
	Router.OnMatch("gif", dgrouter.NewRegexMatcher("^(?i)gif$"), cmd.Gif)
	Router.OnMatch("owofy", dgrouter.NewRegexMatcher("^(?i)owofy$"), cmd.Owofy)
	Router.OnMatch("8ball", dgrouter.NewRegexMatcher("^(?i)(8ball|8b|eightball)$"), cmd.EightBall)
	Router.OnMatch("kaomoji", dgrouter.NewRegexMatcher("^(?i)kaomoji$"), cmd.Kaomoji)

	// Weeb
	Router.OnMatch("catgirl", dgrouter.NewRegexMatcher("^(?i)catgirl$"), cmd.CatGirl)
	Router.OnMatch("headpat", dgrouter.NewRegexMatcher("^(?i)(headpat|pat)$"), cmd.HeadPat)
	Router.OnMatch("kiss", dgrouter.NewRegexMatcher("^(?i)kiss$"), cmd.Kiss)
	Router.OnMatch("tickle", dgrouter.NewRegexMatcher("^(?i)tickle$"), cmd.Tickle)
	Router.OnMatch("feed", dgrouter.NewRegexMatcher("^(?i)(feed|food|eat)$"), cmd.Feed)
	Router.OnMatch("slap", dgrouter.NewRegexMatcher("^(?i)slap$"), cmd.Slap)
	Router.OnMatch("cuddle", dgrouter.NewRegexMatcher("^(?i)cuddle$"), cmd.Cuddle)
	Router.OnMatch("hug", dgrouter.NewRegexMatcher("^(?i)hug$"), cmd.Hug)
	Router.OnMatch("smug", dgrouter.NewRegexMatcher("^(?i)smug$"), cmd.Smug)
	Router.OnMatch("baka", dgrouter.NewRegexMatcher("^(?i)(baka|idiot)$"), cmd.Baka)

	// Music
	// Router.OnMatch("add", dgrouter.NewRegexMatcher("^(?i)add$"), cmd.AddMusic)
	// Router.OnMatch("clear", dgrouter.NewRegexMatcher("^(?i)clear$"), cmd.ClearCommand)
	// Router.OnMatch("current", dgrouter.NewRegexMatcher("^(?i)current$"), cmd.CurrentCommand)
	// Router.OnMatch("join", dgrouter.NewRegexMatcher("^(?i)join$"), cmd.JoinCommand)
	// Router.OnMatch("leave", dgrouter.NewRegexMatcher("^(?i)(leave|disconnect)$"), cmd.LeaveCommand)
	Router.OnMatch("pause", dgrouter.NewRegexMatcher("^(?i)pause$"), cmd.PauseCommand)
	Router.OnMatch("play", dgrouter.NewRegexMatcher("^(?i)play$"), cmd.PlayCommand)
	Router.OnMatch("resume", dgrouter.NewRegexMatcher("^(?i)(resume|continue)$"), cmd.PlayCommand)
	// Router.OnMatch("queue", dgrouter.NewRegexMatcher("^(?i)queue$"), cmd.QueueCommand)
	// Router.OnMatch("shuffle", dgrouter.NewRegexMatcher("^(?i)shuffle$"), cmd.ShuffleCommand)
	// Router.OnMatch("skip", dgrouter.NewRegexMatcher("^(?i)skip$"), cmd.SkipCommand)
	Router.OnMatch("stop", dgrouter.NewRegexMatcher("^(?i)stop$"), cmd.StopCommand)

	// Help
	Router.Default = Router.OnMatch("help", dgrouter.NewRegexMatcher("^(?i)(help|h)$"), func(ctx *exrouter.Context) {
		msg := "Please visit __https://darenliang.github.io/MikuBot-Docs__ for help on all the commands."

		// If DMs
		if ctx.Msg.GuildID != "" {
			msg = fmt.Sprintf("The current server prefix is %s\n", framework.PDB.GetPrefix(ctx.Msg.GuildID)) + msg
		} else {
			msg = fmt.Sprintf("The DM prefix is %s\n", config.Prefix) + msg
		}
		_, _ = ctx.Reply(msg)
	}).Cat("Help")

	// Query database on ready
	Session.AddHandlerOnce(func(_ *discordgo.Session, ready *discordgo.Ready) {
		// Query databases to temp
		framework.PDB.SetGuilds()
		framework.MQDB.SetScores()
		framework.GBD.SetAlbums()

		// Set BotID
		config.BotID = Session.State.User.ID

		// Set music
		cmd.AudioInit()

		// Music sessions and youtube
		music.MusicSessions = music.NewSessionManager()

		// Load cache and check for new guilds
		cache := framework.PDB.GetGuilds()
		for _, guild := range ready.Guilds {
			readyGuilds[guild.ID] = true
			if cache[guild.ID] == "" {
				framework.PDB.CreateGuild(guild.ID, config.Prefix)
			}
		}

		// Set ready
		status.setReady()
	})

	// Add guild on guild add
	Session.AddHandler(func(_ *discordgo.Session, create *discordgo.GuildCreate) {
		if status.getReady() && !framework.PDB.CheckGuild(create.ID) {
			framework.PDB.CreateGuild(create.ID, config.Prefix)
		}
	})

	// Remove guild on guild remote
	Session.AddHandler(func(_ *discordgo.Session, delete *discordgo.GuildDelete) {
		if status.getReady() && framework.PDB.CheckGuild(delete.ID) && !delete.Unavailable {
			framework.PDB.RemoveGuild(delete.ID)
		}
	})

	// Handle incoming messages
	Session.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		prefix := config.Prefix
		if m.GuildID != "" {
			prefix = framework.PDB.GetPrefix(m.GuildID)
		}
		_ = Router.FindAndExecute(Session, prefix, Session.State.User.ID, m.Message)
	})

	// Handle reaction add for gif command
	Session.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		// If DM
		if reaction.GuildID == "" {
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

	Session.AddHandler(func(session *discordgo.Session, event *discordgo.VoiceServerUpdate) {
		var err error

		vsu := gavalink.VoiceServerUpdate{
			Endpoint: event.Endpoint,
			GuildID:  event.GuildID,
			Token:    event.Token,
		}

		if cmd.AudioPlayers[event.GuildID], err = cmd.AudioLavalink.GetPlayer(event.GuildID); err == nil {
			if err = cmd.AudioPlayers[event.GuildID].Forward(session.State.SessionID, vsu); err != nil {
				log.Printf("handler: voice server update: %s", err)
				return
			}
		}

		cmd.AudioNode, err = cmd.AudioLavalink.BestNode()
		if err != nil {
			log.Printf("handler: failed to find node: %s", err)
			return
		}

		if cmd.AudioPlayers[event.GuildID], err = cmd.AudioNode.CreatePlayer(event.GuildID, session.State.SessionID, vsu, new(gavalink.DummyEventHandler)); err != nil {
			log.Printf("handler: failed to create player: %s", err)
			return
		}
	})
}
