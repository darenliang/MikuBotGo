package config

import (
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/framework"
	"time"
)

// Includes config variables for bot
const BotInfo = "MikuBotGo v0.1.0"
const EmbedColor = 0x2e98a6
const Prefix = ";"
const Timeout = 60

var StartTime time.Time
var Openings framework.Openings3
var OpeningsMap = make(map[string]OpeningEntry)

type OpeningEntry struct {
	Name  string
	Embed *discordgo.MessageEmbed
}

func init() {
	// Set start time
	StartTime = time.Now()

	// Setup openings
	Openings = framework.GetOpenings3()
}
