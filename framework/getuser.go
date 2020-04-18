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
	if len(msg.GuildID) == 0 {
		// Recipients don't include the bot itself
		if matchUser(ctx.Ses.State.User, user) {
			return ctx.Ses.State.User
		}

		dm, err := ctx.Ses.Channel(ctx.Msg.ChannelID)

		if err != nil {
			return nil
		}

		for _, u := range dm.Recipients {
			if matchUser(u, user) {
				return u
			}
		}

		return nil
	}

	// Guild channels
	g, err := ctx.Ses.Guild(msg.GuildID)

	if err != nil {
		return nil
	}

	for _, member := range g.Members {
		u := member.User
		if matchUser(u, user) {
			return u
		}
	}

	return nil
}

func matchUser(u *discordgo.User, uString string) bool {
	if u.Username == uString || u.String() == uString || u.ID == uString || getMentionString(u) == uString {
		return true
	}
	return false
}

func getMentionString(u *discordgo.User) string {
	return "<@!" + u.ID + ">"
}
