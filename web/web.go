package web

import (
	"github.com/bwmarrin/discordgo"
	"github.com/cam-per/licob/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	engine  *gin.Engine
	discord *discordgo.Session
)

func Init(s *discordgo.Session) {
	discord = s

	gin.SetMode(gin.DebugMode)
	engine = gin.Default()
	engine.Use(cors.Default())
	engine.Static("/assets", "../../interface/dist/assets")
	engine.LoadHTMLFiles("../../interface/dist/index.html")

	configureRouter()
}

func Run(addr ...string) error {
	if len(addr) == 0 {
		addr = append(addr, ":8080")
	}
	return engine.Run(addr...)
}

func RunAsync() (errCh chan error) {
	errCh = make(chan error)

	go func() {
		errCh <- Run(":" + config.Port)
	}()

	return
}
