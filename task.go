package main

import (
	"fmt"
	"github.com/DiscordBotList/go-dbl"
	"log"
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

// ScheduleGC forces GC to reap the heap every 1 hour
func ScheduleGC() {
	var memRuntime runtime.MemStats
	for {
		time.Sleep(time.Hour)
		runtime.ReadMemStats(&memRuntime)
		log.Printf("Pre-GC heap mem usage: %v MiB", memRuntime.HeapAlloc/1024/1024)
		runtime.GC()
		log.Print("Undergoing GC...")
		runtime.ReadMemStats(&memRuntime)
		log.Printf("Post-GC heap mem usage: %v MiB", memRuntime.HeapAlloc/1024/1024)
	}
}
