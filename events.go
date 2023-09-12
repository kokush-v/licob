package licob

import (
	"github.com/bwmarrin/discordgo"
)

const (
	EventDrop   = "sakura.drop"
	EventPick   = "sakura.pick"
	EventPicked = "sakura.picked"
	EventWin    = "sakura.win"
)

type DropData struct {
	M          *discordgo.Message
	Currency   int
	CaptchaURL string
}

type PickData struct {
	M    *discordgo.Message
	Code string
}

type PickedData struct {
	M *discordgo.Message
}

type WinData struct {
	M        *discordgo.Message
	WinnerID string
}
