package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/configs"
	"log"
	"os"
)

// Create session
var Session, _ = discordgo.New()

// Read in all configuration options from both environment variables and command line arguments.
func init() {
	// Discord Authentication Token
	Session.Token = "Bot " + os.Getenv("DISCORD_TOKEN")
	if Session.Token == "" {
		flag.StringVar(&Session.Token, "t", "", "Discord Authentication Token")
	}
}

func main() {
	// Declare any variables needed later.
	var err error

	// Print bot info
	fmt.Println(configs.BotInfo)

	// Parse command line arguments
	flag.Parse()

	// Open a websocket connection to Discord
	err = Session.Open()
	if err != nil {
		log.Printf("Error opening connection to Discord, %s\n", err)
		return
	}

	// Run and keep open
	log.Printf("Now running.")
	<-make(chan struct{})
}
