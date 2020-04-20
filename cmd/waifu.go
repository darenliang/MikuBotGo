package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Waifu command
func Waifu(ctx *exrouter.Context) {
	imageID := rand.Int() % 100001

	resp, err := http.Get(fmt.Sprintf(
		"https://www.thiswaifudoesnotexist.net/example-%d.jpg", imageID))

	if err != nil {
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Cannot generate an image.")
		return
	}

	_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "Here's your waifu.")
	_, _ = ctx.Ses.ChannelFileSend(ctx.Msg.ChannelID, "waifu.jpg", resp.Body)
	_ = resp.Body.Close()
}
