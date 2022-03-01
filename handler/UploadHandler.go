package handler

import (
	"GoFileService/config"
	"github.com/gin-gonic/gin"
	"log"
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
	fileHeaders := form.File["files"]
	var accessPaths = make([]string, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		// 先直接保存到本地
		err = c.SaveUploadedFile(fileHeader, dirPath+"/"+fileHeader.Filename)
		if err != nil {
			handleErr(c, 500, err)
			return
		}

		// 是图片且需要经过压缩处理
		if IsImage(fileHeader.Filename) && fileHeader.Size > (1<<20) {
			go func() {
				err := HandleImageFile(dirPath, fileHeader.Filename)
				if err != nil {
					log.Printf("压缩图片出现异常：%+v", err)
				}
			}()
			accessPaths = append(accessPaths, accessPath+"/"+"_"+fileHeader.Filename) // 被压缩的图片添加访问路劲，"_"作前缀
		} else {
			accessPaths = append(accessPaths, accessPath+"/"+fileHeader.Filename) // 添加访问路劲
		}
	}
	c.JSON(200, Result{200, "ok", accessPaths})
}

func handleErr(c *gin.Context, code int, err error) {
	c.JSON(code, Result{code, err.Error(), nil})
}
