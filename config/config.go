package config

import (
	"github.com/darenliang/MikuBotGo/framework"
	"time"
)

type OpeningsEntry struct {
	Id      int
	Answers []string
	Source  string
}

// Includes config variables for bot
const BotInfo = "MikuBotGo v0.0.0-alpha"
const EmbedColor = 0x2e98a6
const Prefix = ";"
const Timeout = 60
const TimeoutMsg = "Your command has timed out"

var StartTime time.Time
var Openings framework.Openings
var OpeningsMap = make(map[string]OpeningsEntry)

func init() {
	// Set start time
	StartTime = time.Now()

	// Setup openings
	Openings = framework.GetOpenings()
}
