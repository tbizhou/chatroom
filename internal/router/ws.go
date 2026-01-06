package router

import (
	"github.com/chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

func initWebsocketRouter(r *gin.Engine) {
	r.GET("/ws/:id", service.HandleWebSocket)
}
