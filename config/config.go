package config

import (
	"github.com/darenliang/MikuBotGo/framework"
	"time"
)

type OpeningsEntry struct {
	Answers []string
	Source  string
}

// Includes config variables for bot
const BotInfo = "MikuBotGo v0.1.0"
const EmbedColor = 0x2e98a6
const Prefix = ";"
const Timeout = 60

var StartTime time.Time
var Openings framework.Openings2
var OpeningsMap = make(map[string]OpeningsEntry)

func init() {
	// Set start time
	StartTime = time.Now()

	// Setup openings
	Openings = framework.GetOpenings2()
}
