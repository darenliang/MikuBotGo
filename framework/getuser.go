package framework

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func Getuser(ctx *exrouter.Context) *discordgo.User {
	msg := ctx.Msg

	// Get full string after command invoke
	user := strings.TrimSpace(ctx.Args.After(1))

	// If no arguments
	if len(user) == 0 {
		return msg.Author
	}

	// Direct messages
	if msg.GuildID == "" {
		// Recipients don't include the bot itself
		if matchUser(ctx.Ses.State.User, user) {
			return ctx.Ses.State.User
		}

		dm, err := ctx.Ses.State.Channel(ctx.Msg.ChannelID)
		if err != nil {
			dm, err = ctx.Ses.Channel(ctx.Msg.ChannelID)
			if err != nil {
				return nil
			}
			ctx.Ses.State.ChannelAdd(dm)
		}

		for _, u := range dm.Recipients {
			if matchUser(u, user) {
				return u
			}
		}

		return nil
	}

	// Guilds
	guild, err := ctx.Ses.State.Guild(msg.GuildID)
	if err != nil {
		guild, err = ctx.Ses.Guild(msg.GuildID)
		if err != nil {
			return nil
		}
		ctx.Ses.State.GuildAdd(guild)
	}

	for _, member := range guild.Members {
		if member.Nick == user || matchUser(member.User, user) {
			return member.User
		}
	}

	return nil
}

func matchUser(u *discordgo.User, uString string) bool {
	return u.Username == uString || u.String() == uString || u.ID == uString || getMentionString(u) == uString
}

func getMentionString(u *discordgo.User) string {
	return "<@!" + u.ID + ">"
}
