package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/framework"
)

func Pfp(ctx *exrouter.Context) {

	target := framework.Getuser(ctx)
	if target != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, target.AvatarURL("1024"))
	} else {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred.")
	}

}
