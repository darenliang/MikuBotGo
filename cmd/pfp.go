package cmd

import (
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/darenliang/MikuBotGo/framework"
)

func Pfp(ctx *exrouter.Context) {

	target := framework.Getuser(ctx)
	if target != nil {
		ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, target.AvatarURL("1024"))
	} else {
		ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occured.")
	}

}
