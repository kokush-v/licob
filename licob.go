package licob

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

func Auth() (*discordgo.Session, error) {
	readyCh := make(chan struct{})
	onReady := func(s *discordgo.Session, ready *discordgo.Ready) {
		readyCh <- struct{}{}
	}

	discord := discordgo.New()
	discord.ShouldReconnectOnError = false
	discord.ShouldRetryOnRateLimit = false
	discord.Identify.Intents = discordgo.IntentsAll

	discord.AddHandlerOnce(onReady)
	discord.AddHandler(handleDrop)
	discord.AddHandler(handleWin)
	discord.AddHandler(handlePick)
	discord.AddHandler(handlePicked)
	var (
		email    = os.Getenv("LOGIN")
		password = os.Getenv("PASSWORD")
	)

	ticked, err := discord.Login(email, password)
	if err != nil {
		return nil, err
	}

	if ticked != "" {
		var code string
		fmt.Print("2FA eneblad. Ented auth code: ")
		fmt.Scan(&code)

		err = discord.TOTP(ticked, code)
		if err != nil {
			return nil, err
		}
	}

	err = discord.Open()
	if err != nil {
		return nil, err
	}

	select {
	case <-time.After(25 * time.Second):
		return nil, errors.New("open timeout")
	case <-readyCh:
	}

	return discord, nil
}
