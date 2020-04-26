package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"regexp"
	"strings"
)

// Owofy command. Ported from kyostra/owofy
func Owofy(ctx *exrouter.Context) {
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Please provide a message to owofy.")
		return
	}

	var re = regexp.MustCompile("[rl]")
	query = re.ReplaceAllString(query, "w")
	re = regexp.MustCompile("[RL]")
	query = re.ReplaceAllString(query, "W")
	re = regexp.MustCompile("n([aeiouAEIOU])")
	query = re.ReplaceAllString(query, "ny$1")
	re = regexp.MustCompile("N([aeiouAEIOU])")
	query = re.ReplaceAllString(query, "Ny$1")
	query = strings.ReplaceAll(query, "ove", "uv")

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, query)
}
