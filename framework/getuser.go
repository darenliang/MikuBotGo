package framework

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
)

func Getuser(ctx *exrouter.Context) *discordgo.User {

	msg := ctx.Msg

	if len(msg.Mentions) > 0 && msg.Mentions[0].ID != ctx.Ses.State.User.ID {
		return msg.Mentions[0]
	}

	if len(ctx.Args) < 2 {
		return msg.Author
	}

	g, err := ctx.Ses.Guild(msg.GuildID)

	if err != nil {
		return nil
	}

	user := ctx.Args[1]

	for _, member := range g.Members {
		u := member.User
		if u.Username == user || u.String() == user || u.ID == user {
			return member.User
		}
	}

	if len(ctx.Args) > 1 {
		return nil
	}

	return msg.Author

}
