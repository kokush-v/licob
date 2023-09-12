package licob

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/cam-per/licob/config"
	"github.com/gookit/event"
)

func handleDrop(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.DropBotID {
		return
	}

	const format = "%d случайных 🌸 появились! Напишите `.pick и код с картинки`, чтобы собрать их."
	var currency int
	n, err := fmt.Sscanf(m.Content, format, &currency)
	if err != nil || n != 1 {
		return
	}

	if len(m.Attachments) != 1 {
		return
	}

	event.Fire(EventDrop, event.M{
		"data": &DropData{
			M:          m.Message,
			Currency:   currency,
			CaptchaURL: m.Attachments[0].URL,
		},
	})
}

func handlePicked(s *discordgo.Session, m *discordgo.MessageDelete) {
	event.Fire(EventPicked, event.M{
		"data": &PickedData{
			M: m.Message,
		},
	})

}

func handlePick(s *discordgo.Session, m *discordgo.MessageCreate) {
	var code string

	n, err := fmt.Sscanf(m.Content, ".pick %s", &code)
	if err != nil || n != 1 {
		return
	}
	if len([]rune(code)) != 4 {
		return
	}

	event.Fire(EventPick, event.M{
		"data": &PickData{
			M:    m.Message,
			Code: code,
		},
	})
}

func handleWin(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != config.DropBotID {
		return
	}
	if len(m.Embeds) != 1 {
		return
	}

	var id, currency uint64

	n, err := fmt.Sscanf(m.Embeds[0].Description, "**<@!%d>** собрал %d🌸", &id, &currency)
	if err != nil || n != 2 {
		return
	}

	event.Fire(EventWin, event.M{
		"data": &WinData{
			M:        m.Message,
			WinnerID: strconv.FormatUint(id, 10),
		},
	})
}
