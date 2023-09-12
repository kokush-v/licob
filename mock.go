package licob

import "github.com/bwmarrin/discordgo"

func mockDrop() *discordgo.MessageCreate {
	const format = "845 случайных 🌸 появились! Напишите `.pick и код с картинки`, чтобы собрать их."

	result := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ID:        "mock_drop_id",
			ChannelID: "channel_mock",
			Content:   "845 случайных 🌸 появились! Напишите `.pick и код с картинки`, чтобы собрать их.",
		},
	}
	return result
}
