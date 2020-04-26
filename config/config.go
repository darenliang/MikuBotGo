package config

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/music"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

// Includes config variables for bot
const (
	BotInfo              = "MikuBot v1.3.0"
	EmbedColor           = 0x2e98a6
	Prefix               = "!"
	Timeout              = 60
	Timer                = "\xe2\x8f\xb2\xef\xb8\x8f"
	ImgurEndpoint        = "https://api.imgur.com/3"
	MaxImgurByteSize     = 1000 * 1000 * 10
	ClarifaiNSFWEndpoint = "https://api.clarifai.com/v2/models/e9576d86d2004ed1a38ba0cf39ecb4b1/outputs"
)

var (
	StartTime    time.Time
	OpeningsData Openings
	OpeningsMap  = sync.Map{}

	// TODO: Make thread safe
	MusicSessions *music.SessionManager

	TraceMoeBase  string
	ImgurToken    string
	ImgurUsername string
	ClarifaiToken string
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
	ClarifaiToken = os.Getenv("CLARIFAI_TOKEN")
}
