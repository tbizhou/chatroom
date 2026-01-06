package router

import (
	"github.com/chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

func initMessageRouter(r *gin.Engine) {
	message := r.Group("/api/message")
	{
		message.GET("/:id", service.GetMessages)
	}
}
