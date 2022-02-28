package handler

import (
	"GoFileService/config"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

// Result 自定义结果返回
type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var myConfig *config.Config = nil

func RegisterUploadHandler(engine *gin.Engine, myConf *config.Config) {
	myConfig = myConf
	engine.POST("/file/upload", uploadHandle)
}

func uploadHandle(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		handleErr(c, 400, err)
		return
	}
	// 获取文件保存路径和访问路径
	var relativePath = c.PostForm("relativePath")
	var fileType = c.PostForm("type")

	var dirPath = myConfig.DirRootPath + "/" + fileType
	var accessPath = myConfig.DirAccessPath + "/" + fileType
	if strings.HasPrefix(relativePath, "/") {
		dirPath = dirPath + relativePath
		accessPath = accessPath + relativePath
	} else {
		dirPath = dirPath + "/" + relativePath
		accessPath = accessPath + "/" + relativePath
	}
	err = os.MkdirAll(dirPath, os.ModePerm) // 创建目录
	if err != nil {
		handleErr(c, 500, err)
		return
	}

	// 保存文件到本地
	files := form.File["files"]
	var accessPaths = make([]string, 0, len(files))
	for _, file := range files {
		err = c.SaveUploadedFile(file, dirPath+"/"+file.Filename)
		accessPaths = append(accessPaths, accessPath+"/"+file.Filename) // 添加访问路劲
		if err != nil {
			handleErr(c, 500, err)
			return
		}
	}

	c.JSON(200, Result{200, "ok", accessPaths})
}

func handleErr(c *gin.Context, code int, err error) {
	c.JSON(code, Result{code, err.Error(), nil})
}
