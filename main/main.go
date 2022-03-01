package main

import (
	"GoFileService/config"
	"GoFileService/handler"
	"GoFileService/router"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

func main() {
	serverConfig := config.GetServerConfig()
	engine := gin.Default()
	router.RegisterRouter(engine)       // 注册路由：首页
	handler.RegisterFileHandler(engine) // 注册路由：文件处理器

	err := engine.Run(":" + strconv.Itoa(serverConfig.Port))
	if err != nil {
		log.Fatalln(err)
	}

}
