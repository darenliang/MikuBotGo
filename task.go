package main

import (
	"fmt"
	"time"
)

// UpdatePresence updates the presence every 15 minutes
func UpdatePresence() {
	for {
		_ = Session.UpdateStatus(0, fmt.Sprintf("@%s help", Session.State.User.Username))
		time.Sleep(time.Minute * 15)
	}
}
