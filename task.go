package main

import (
	"time"
)

// UpdatePresence updates the presence every 15 minutes
func UpdatePresence() {
	for {
		_ = Session.UpdateStatus(0, "@MikuBot help")
		time.Sleep(time.Minute * 15)
	}
}
