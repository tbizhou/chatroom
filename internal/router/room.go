package router

import (
	"github.com/chatroom/internal/service"
	"github.com/gin-gonic/gin"
)

func initRoomRouter(r *gin.Engine) {
	room := r.Group("/api/room")
	{
		room.POST("", service.CreateRoom) //creat chat room
		room.GET("", service.GetRooms)    //get chat rooms
	}
}
