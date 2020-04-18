package main

import (
	"fmt"
	"github.com/darenliang/MikuBotGo/config"
	"time"
)

// UpdatePresence updates the presence every hour
func UpdatePresence() {
	for {
		_ = Session.UpdateStatus(0, fmt.Sprintf("%shelp", config.Prefix))
		time.Sleep(time.Hour)
	}
}
