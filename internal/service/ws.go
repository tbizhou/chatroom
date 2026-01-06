package service

import (
	"log"
	"net/http"

	"github.com/chatroom/internal/data"
	"github.com/chatroom/pkg/storage"
	"github.com/chatroom/pkg/websocket"
	"github.com/gin-gonic/gin"
	ws "github.com/gorilla/websocket"
)

func init() {
	// 初始化Hub的消息存储
	hub := websocket.GetHub()
	hub.SetMessageStore(data.NewRedisClient(storage.GetRedis()))
}

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，生产环境应该限制
	},
}

func HandleWebSocket(c *gin.Context) {
	roomID := c.Param("id")
	userID := c.Query("user_id")
	if userID == "" {
		userID = "anonymous"
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	hub := websocket.GetHub()
	room := hub.GetRoom(roomID)
	client := websocket.NewClient(hub, conn, room, userID)

	room.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
