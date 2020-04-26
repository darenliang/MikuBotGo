package main

import (
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	"os"
	"time"
)

var (
	dblToken  string
	dblClient *dbl.Client
)

func init() {
	dblToken = os.Getenv("DBL_TOKEN")

	if dblToken != "" {
		dblClient, _ = dbl.NewClient(dblToken)
	}
}

// UpdatePresence updates the presence every 15 minutes
func UpdatePresence() {
	for {
		if dblToken != "" {
			_ = dblClient.PostBotStats(Session.State.User.ID, &dbl.BotStatsPayload{
				Shards: []int{len(Session.State.Guilds)},
			})
		}
		_ = Session.UpdateStatus(0, fmt.Sprintf("@%s help", Session.State.User.Username))
		time.Sleep(time.Minute * 15)
	}
}
