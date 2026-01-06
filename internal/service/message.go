package service

import (
	"github.com/chatroom/utils"
	"github.com/gin-gonic/gin"
)

func GetMessages(c *gin.Context) {
	roomId := c.Param("id")
	result, err := rdb.GetRoomMessages(c, roomId, 50)
	if err != nil {
		utils.InternalServerError(c, "get chat room message failed")
	}
	utils.SuccessWithData(c, result)
}
