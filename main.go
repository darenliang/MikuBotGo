package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"log"
	"os"
)

// Create session
var Session, _ = discordgo.New()

// Initialize flag token
var FlagToken = flag.String("t", "", "Discord Authentication Token")

// Read in all configuration options from both environment variables and command line arguments.
func init() {
	// Discord Authentication Token
	EnvToken := os.Getenv("DISCORD_TOKEN")
	if EnvToken != "" {
		Session.Token = EnvToken
	} else {
		// Parse command line arguments
		flag.Parse()
		Session.Token = *FlagToken
	}
	Session.Token = "Bot " + Session.Token
}

func main() {
	// Declare any variables needed later.
	var err error

	// Print bot info
	fmt.Println(config.BotInfo)

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("Error opening connection to Discord, %s\n", err)
		return
	}

	go UpdatePresence()
	go ScheduleGC()

	// Run and keep open
	log.Printf("Now running.")
	<-make(chan struct{})
}
