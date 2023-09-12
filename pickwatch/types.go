package pickwatch

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/cam-per/licob/utils"
)

type Channel struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type User struct {
	ID   string `json:"id"`
	Nick string `json:"nick"`
}

type Drop struct {
	t time.Time

	ID        string `json:"id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Currency  int    `json:"currency,omitempty"`
	Since     string `json:"since,omitempty"`

	Picks  []*Pick `json:"picks,omitempty"`
	Winner *User   `json:"winner,omitempty"`
}

type Pick struct {
	t time.Time

	ID        string `json:"id,omitempty"`
	DropID    string `json:"drop_id,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Code      string `json:"code,omitempty"`
	Ok        bool   `json:"ok,omitempty"`
	Since     string `json:"since,omitempty"`
	User      *User  `json:"user,omitempty"`
}

type EventWin Drop

func newChannel(channel *discordgo.Channel) *Channel {
	return &Channel{
		ID:   channel.ID,
		Name: channel.Name,
	}
}

func newUser(user *discordgo.User) *User {
	return &User{
		ID:   user.ID,
		Nick: user.String(),
	}
}

func newDrop(message *discordgo.Message) *Drop {
	t := utils.SnowflakeTimestamp(message.ID)
	return &Drop{
		t: t,

		ID:        message.ID,
		Timestamp: t.Unix(),
	}
}

func newPick(drop *Drop, message *discordgo.Message) *Pick {
	t := utils.SnowflakeTimestamp(message.ID)
	since := t.Sub(drop.t)
	return &Pick{
		t: t,

		ID:        message.ID,
		DropID:    drop.ID,
		Timestamp: t.Unix(),
		Since:     since.Round(time.Millisecond).String(),
		User:      newUser(message.Author),
	}
}
