package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func configureRouter() {
	engine.Use(gin.ErrorLoggerT(gin.ErrorTypeAny))

	engine.GET("/api/channels/:channel_id", getChannel)
	engine.GET("/ws", ws)
	engine.GET("/channels/:channel_id", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	})
}
