package main

import (
	"GoFileService/config"
	"GoFileService/handler"
	"GoFileService/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	myConfig := config.LoadConfig() // 加载配置文件
	engine := gin.Default()
	router.RegisterRouter(engine)                   // 注册首页
	handler.RegisterUploadHandler(engine, myConfig) // 注册上传处理函数

	err := engine.Run(":" + myConfig.Port)
	if err != nil {
		fmt.Println(err)
		return
	}
}
