package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getChannel(c *gin.Context) {
	id := c.Param("channel_id")

	ch, err := discord.Channel(id)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	me, err := discord.User("@me")
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"channel": channel{
			ID:   ch.ID,
			Name: ch.Name,
		},
		"user": user{
			ID:     me.ID,
			Nick:   me.String(),
			Avatar: me.AvatarURL(""),
		},
	})
}
