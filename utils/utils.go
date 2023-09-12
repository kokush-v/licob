package utils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func AlloweChannels(channel string, allowed []string) bool {
	for _, id := range allowed {
		if channel == id {
			return true
		}
	}
	return false
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func SnowflakeTimestamp(id string) time.Time {
	resp, err := discordgo.SnowflakeTimestamp(id)
	if err != nil {
		resp = time.Now()
	}
	return resp
}
