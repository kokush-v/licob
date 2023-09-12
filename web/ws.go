package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/cam-per/licob"
	"github.com/cam-per/licob/picker"
	"github.com/cam-per/licob/pickwatch"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func reader(conn *websocket.Conn, discord *discordgo.Session, c *gin.Context, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	for {
		mt, mc, err := conn.ReadMessage()
		if err != nil {
			return
		}

		fmt.Println(mt, string(mc))

		var data map[string]interface{}

		err = json.Unmarshal([]byte(string(mc)), &data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if data["t"] == "sakura.send" {
			channelID := c.Query("channel_id")
			discord.ChannelMessageSend(channelID, fmt.Sprintf(".pick %s", data["d"]))
		}
	}
}

func writer(conn *websocket.Conn, evch <-chan any, done chan struct{}) {
	defer func() {
		done <- struct{}{}
	}()

	for data := range evch {
		var t string

		switch data.(type) {
		default:
			continue

		case *pickwatch.Drop:
			t = licob.EventDrop
		case *pickwatch.Pick:
			t = licob.EventPick
		case *pickwatch.EventWin:
			t = licob.EventWin

		case *picker.EventCaptcha:
			t = "sakura.captcha"
		case *picker.EventPicked:
			t = "sakura.picked"
		}

		data := map[string]any{
			"t": t,
			"d": data,
		}
		conn.WriteJSON(data)
	}
}

func ws(c *gin.Context) {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	channelID := c.Query("channel_id")

	log.Println("socket connected", conn.RemoteAddr())

	evch := make(chan any)
	defer close(evch)

	pickwatcher := pickwatch.NewPickwatch(&pickwatch.Opt{
		ChannelID: channelID,
		Events:    evch,
	})
	picker := picker.NewPicker(&picker.Opt{
		ChannelID: channelID,
		Events:    evch,
	})

	pickwatcher.Listen()
	picker.Listen()
	defer pickwatcher.Close()
	defer picker.Close()

	done := make(chan struct{})
	go reader(conn, discord, c, done)
	go writer(conn, evch, done)
	<-done
}
