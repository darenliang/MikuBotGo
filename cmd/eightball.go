package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/framework"
	"math/rand"
	"strings"
)

var eightBallChoices = [...]string{"It is certain.", "It is decidedly so.", "Without a doubt.", "Yes - definitely.",
	"You may rely on it.", "As I see it, yes.", "Most likely.", "Outlook good.",
	"Yes.", "Signs point to yes.", "Reply hazy, try again.", "Ask again later.",
	"Better not tell you now.", "Cannot predict now.", "Concentrate and ask again.",
	"Don't count on it.", "My reply is no.", "My sources say no.", "Outlook not so good.",
	"Very doubtful."}

// EightBall command.
func EightBall(ctx *exrouter.Context) {
	prefix := framework.PDB.GetPrefix(ctx.Msg.GuildID)
	query := strings.TrimSpace(ctx.Args.After(1))

	if len(query) == 0 {
		ctx.Reply(fmt.Sprintf("Usage: `%s8ball <question>`", prefix))
		return
	}

	ctx.Reply(eightBallChoices[rand.Intn(len(eightBallChoices))])
}
