package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRouter(engine *gin.Engine) {
	engine.LoadHTMLGlob("./resources/*")
	engine.GET("/file/index.html", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	engine.GET("/file/index", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/file/index.html")
	})
	engine.GET("/file/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/file/index.html")
	})
}
