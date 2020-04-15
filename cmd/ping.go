package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
)

// Ping command
// This will be removed soon
func Ping(ctx *exrouter.Context) {
	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "pong")
}
