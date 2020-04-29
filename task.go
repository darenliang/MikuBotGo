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

// ReportMem reports mem usage every hour
func ReportMem() {
	var memRuntime runtime.MemStats
	for {
		time.Sleep(time.Minute * 30)
		runtime.ReadMemStats(&memRuntime)
		log.Printf("Heap mem usage: %v MiB", memRuntime.HeapAlloc/1024/1024)
		log.Printf("Sys mem usage: %v MiB", memRuntime.Sys/1024/1024)
		log.Printf("Number of goroutines: %d", runtime.NumGoroutine())
	}
}
