package main

import (
	"GoFileService/config"
	"GoFileService/docs"
	"GoFileService/handler"
	"GoFileService/router"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"strconv"
)

// @title fzk API
// @version 1.0
// @description 文件管理服务API文档
// @BasePath /filemanage 基础路劲
// @host localhost:23456
// @query.collection.format multi
func main() {
	serverConfig := config.GetServerConfig()
	r := gin.Default()
	router.RegisterRouter(r) // 注册路由：首页

	// 注册swagger
	docs.SwaggerInfo.BasePath = "/filemanage"
	handler.RegisterFileHandler(r) // 注册路由：文件处理器
	r.GET("/filemanage/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err := r.Run(":" + strconv.Itoa(serverConfig.Port))
	if err != nil {
		log.Fatalln(err)
	}
}
