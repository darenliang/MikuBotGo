package configs

import "time"

const BotInfo = "MikuBotGo v0.0.0-alpha"
const EmbedColor = 0x2e98a6
const Prefix = "?"
const Timeout = 60
const TimeoutMsg = "Your command has timed out"

var StartTime time.Time

func init() {
	// Set start time
	StartTime = time.Now()
}
