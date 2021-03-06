package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/framework"
	"regexp"
	"strings"
)

// Owofy command. Ported from kyostra/owofy
func Owofy(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf(":information_source: Usage: `%sowofy <message>`", prefix))
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

	ctx.Reply(query)
}
