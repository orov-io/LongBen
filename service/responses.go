package service

import (
	"github.com/gin-gonic/gin"
	"github.com/orov-io/BlackBart/response"
	"github.com/orov-io/LongBen/models"
)

func sendPong(c *gin.Context, pong models.Pong) {
	response.SendOK(c, pong)

}
