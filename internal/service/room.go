package service

import (
	"github.com/chatroom/internal/data"
	"github.com/chatroom/pkg/storage"
	"github.com/chatroom/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var rdb = data.NewRedisClient(storage.GetRedis())

func GetRooms(c *gin.Context) {
	result := rdb.GetRooms(c)
	utils.SuccessWithData(c, result)
}

func CreateRoom(c *gin.Context) {
	roomId := uuid.New().String()
	err := rdb.CreatRoom(c, roomId)
	if err != nil {
		utils.InternalServerError(c, "create chat room failed")
	}
	utils.SuccessWithData(c, roomId)
}
