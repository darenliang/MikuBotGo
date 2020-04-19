package main

import (
	"time"
)

// UpdatePresence updates the presence every hour
func UpdatePresence() {
	for {
		_ = Session.UpdateStatus(0, "@MikuBotGo help")
		time.Sleep(time.Hour)
	}
}
