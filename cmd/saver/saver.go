package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/cam-per/licob/picker"
	"github.com/cam-per/licob/utils"
)

type Saver struct {
	discord *discordgo.Session
}

func NewSaver(session *discordgo.Session) *Saver {
	saver := &Saver{
		discord: session,
	}
	return saver
}

func (saver *Saver) Start() {
	os.Mkdir("images", os.ModePerm)
	os.Mkdir("out", os.ModePerm)
	saver.discord.AddHandler(saver.gateway)
}

func (saver *Saver) procDrop(m *discordgo.Message) {
	var drop int
	n, err := fmt.Sscanf(m.Content, "%d случайных 🌸 появились! Напишите `.pick и код с картинки`, чтобы собрать их.", &drop)
	if n != 1 || err != nil {
		return
	}
	log.Println(m.Content)

	if len(m.Attachments) == 0 {
		return
	}

	attach := m.Attachments[0]
	if attach == nil {
		return
	}

	log.Println(attach.URL)
	resp, err := http.Get(attach.URL)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	fname := filepath.Join(".", "images", m.ID+".png")
	f, err := os.Create(fname)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, resp.Body)
	f.Seek(0, io.SeekStart)

	captcha, err := picker.NewCaptchaDecoder(f)
	if err != nil {
		log.Println(err)
		return
	}

	err = captcha.Decode()
	if err != nil {
		captcha.SaveDebug(m.ID)
		log.Println(err)
		return
	}

	fmt.Println(captcha.Codes())
	captcha.SaveDebug(m.ID)
}

func (saver *Saver) gateway(s *discordgo.Session, m *discordgo.MessageCreate) {
	allow := utils.AlloweChannels(m.ChannelID, strings.Fields(os.Getenv("DROP_CHANNELS")))

	if !allow || m.Author.ID != os.Getenv("DROP_BOT_ID") {
		return
	}

	saver.procDrop(m.Message)
}
