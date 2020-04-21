package config

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// Includes config variables for bot
const BotInfo = "MikuBotGo v1.0.0"
const EmbedColor = 0x2e98a6
const Prefix = ";"
const Timeout = 60
const Timer = "\xe2\x8f\xb2\xef\xb8\x8f"
const ImgurEndpoint = "https://api.imgur.com/3"
const MaxImgurByteSize = 1000 * 1000 * 10

var (
	StartTime     time.Time
	OpeningsData  Openings
	OpeningsMap   = sync.Map{}
	TraceMoeBase  string
	ImgurToken    string
	ImgurUsername string
)

type OpeningEntry struct {
	Name  string
	Embed *discordgo.MessageEmbed
}

type Openings []struct {
	Name  string `json:"name"`
	Songs []struct {
		Songname string `json:"songname"`
		URL      string `json:"url"`
	} `json:"songs"`
}

// Return openings
func GetOpenings() Openings {
	file, _ := ioutil.ReadFile("data/dataset_filtered.json")
	tmp := Openings{}
	_ = json.Unmarshal(file, &tmp)
	return tmp
}

func init() {
	// Set start time
	StartTime = time.Now()

	// Setup openings
	OpeningsData = GetOpenings()

	TraceMoeKey := os.Getenv("TRACEMOE")
	if TraceMoeKey != "" {
		TraceMoeBase = fmt.Sprintf("https://trace.moe/api/search?token=%s", TraceMoeKey)
	} else {
		TraceMoeBase = "https://trace.moe/api/search"
	}

	ImgurToken = os.Getenv("IMGUR_TOKEN")
	ImgurUsername = os.Getenv("IMGUR_USERNAME")
}
