package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(engine *gin.Engine) {
	engine.LoadHTMLGlob("./resources/*")
	engine.GET("/filemanage/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	engine.GET("/filemanage/list.html", func(c *gin.Context) {
		c.HTML(200, "list.html", nil)
	})
	engine.GET("/filemanage/index", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/filemanage/index.html")
	})
	engine.GET("/filemanage/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/filemanage/index.html")
	})
}
