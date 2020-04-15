package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
)

func Ping(ctx *exrouter.Context) {
	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "pong")
}
