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

// UpdatePresence updates the presence every 1 minute
func UpdatePresence() {
	for {
		if dblToken != "" {
			dblClient.PostBotStats(Session.State.User.ID, &dbl.BotStatsPayload{
				Shards: []int{len(Session.State.Guilds)},
			})
		}
		Session.UpdateStatus(0, fmt.Sprintf("@%s help", Session.State.User.Username))
		time.Sleep(time.Minute)
	}
}
