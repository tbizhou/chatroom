package router

import "github.com/gin-gonic/gin"

func Startup() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// 静态文件服务
	r.Static("/static", "./static")
	r.StaticFile("/", "./static/index.html")

	initWebsocketRouter(r)
	initRoomRouter(r)
	initMessageRouter(r)
	return r
}
