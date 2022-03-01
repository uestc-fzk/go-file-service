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
	myConfig := config.LoadConfig() // 加载配置文件
	engine := gin.Default()
	router.RegisterRouter(engine)                   // 注册首页
	handler.RegisterUploadHandler(engine, myConfig) // 注册上传处理函数

	err := engine.Run(":" + strconv.Itoa(myConfig.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
