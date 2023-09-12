package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	DropChannels []string
	DropBotID    string
	Port         string
)

func Lookup() {
	godotenv.Load()

	DropChannels = strings.Fields(os.Getenv("DROP_CHANNELS"))
	DropBotID = os.Getenv("DROP_BOT_ID")
	Port = os.Getenv("PORT")
}
