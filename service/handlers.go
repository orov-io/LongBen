package service

import (
	"github.com/gin-gonic/gin"
	"github.com/orov-io/LongBen/models"
)

func pong(c *gin.Context) {
	sendPong(c, models.Pong{
		Status:  "OK",
		Message: "pong",
	})
}
