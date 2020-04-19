package config

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/framework"
	"os"
	"time"
)

// Includes config variables for bot
const BotInfo = "MikuBotGo v0.1.0"
const EmbedColor = 0x2e98a6
const Prefix = ";"
const Timeout = 60
const Timer = "\xe2\x8f\xb2\xef\xb8\x8f"

var StartTime time.Time
var Openings framework.Openings3
var OpeningsMap = make(map[string]OpeningEntry)
var TraceMoeBase string

type OpeningEntry struct {
	Name  string
	Embed *discordgo.MessageEmbed
}

func init() {
	// Set start time
	StartTime = time.Now()

	// Setup openings
	Openings = framework.GetOpenings3()

	TraceMoeKey := os.Getenv("TRACEMOE")
	if TraceMoeKey != "" {
		TraceMoeBase = fmt.Sprintf("https://trace.moe/api/search?token=%s", TraceMoeKey)
	} else {
		TraceMoeBase = "https://trace.moe/api/search"
	}
}
