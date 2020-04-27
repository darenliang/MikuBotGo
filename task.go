package main

import (
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	"os"
	"runtime"
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

// ScheduleGC forces a GC every 1 hour
func ScheduleGC() {
	for {
		runtime.GC()
		time.Sleep(time.Hour)
	}
}
