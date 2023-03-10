package handler

import (
	"GoFileService/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"strings"
)

// Result 自定义结果返回
type Result struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"ok"`
	Data interface{} `json:"data"`
}

func RegisterFileHandler(engine *gin.Engine) {
	group := engine.Group("/filemanage")
	{
		group.POST("/upload", UploadHandle)
		group.GET("/fileList", QueryListHandle)
		group.POST("/deleteFileList", DeleteFileListHandler)
	}
}

// UploadHandle
// @Summary      	上传文件
// @Description  	上传普通文件或图片
// @Tags			filemanage
// @Accept       	multipart/form-data
// @Produce      	application/json
// @Param        	relativePath  body string true "相对路径"
// @Param        	type  body string true "文件用途"
// @Success      	200     {object}  handler.Result
// @Router       	/upload [post]
func UploadHandle(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		handleErr(c, 400, err)
		return
	}
	// 1.获取参数得到文件类型和
	var relativePath = c.PostForm("relativePath")
	var fileType = c.PostForm("type")

	// 2.获取文件保存路径和图片保存路径，以及分别的访问路径
	var fileDirPath = config.GetFileConfig().FileDirPath + "/" + fileType
	var fileAccessDirPath = config.GetFileConfig().FileAccessDirPath + "/" + fileType
	var imageDirPath = config.GetImageConfig().ImageDirPath + "/" + fileType
	var imageAccessDirPath = config.GetImageConfig().ImageAccessDirPath + "/" + fileType
	if strings.HasPrefix(relativePath, "/") {
		fileDirPath = fileDirPath + relativePath
		fileAccessDirPath = fileAccessDirPath + relativePath
		imageDirPath = imageDirPath + relativePath
		imageAccessDirPath = imageAccessDirPath + relativePath
	} else {
		fileDirPath = fileDirPath + "/" + relativePath
		fileAccessDirPath = fileAccessDirPath + "/" + relativePath
		imageDirPath = imageDirPath + "/" + relativePath
		imageAccessDirPath = imageAccessDirPath + "/" + relativePath
	}

	// 3. 遍历上传的文件列表
	fileHeaders := form.File["files"]
	var accessPaths = make([]string, 0, len(fileHeaders))
	for _, fileHeader := range fileHeaders {
		// 4. 如果是图片
		if IsImage(fileHeader.Filename) {
			// 4.1 创建目录
			err = os.MkdirAll(imageDirPath, os.ModePerm) // 创建图片保存目录
			if err != nil {
				handleErr(c, 500, err)
				return
			}
			// 4.2 先直接保存到本地
			err = c.SaveUploadedFile(fileHeader, imageDirPath+"/"+fileHeader.Filename)
			if err != nil {
				handleErr(c, 500, err)
				return
			}
			// 4.3 达到压缩限制并且能进行压缩处理如jpeg/jpg和png格式
			if fileHeader.Size > int64(config.GetImageConfig().ToCompressSize) && CanComPress(fileHeader.Filename) {

				// 用协程去缩放图片!!!这里必须将fileHeader指针以参数形式传入进去，不然就会出现闭包现象！！！
				go func(anoFileHeader *multipart.FileHeader) {
					err := CompressImage(imageDirPath+"/"+anoFileHeader.Filename, imageDirPath+"/_"+anoFileHeader.Filename)
					if err != nil {
						log.Printf("压缩图片出现异常：%+v\n", err)
					}
				}(fileHeader)

				accessPaths = append(accessPaths, imageAccessDirPath+"/"+"_"+fileHeader.Filename) // 被压缩的图片添加访问路劲，"_"作前缀
			} else {
				accessPaths = append(accessPaths, imageAccessDirPath+"/"+fileHeader.Filename) // 添加访问路劲
			}
		} else { // 5.其他文件只是少了一个压缩步骤
			// 5.1 创建目录
			err = os.MkdirAll(fileDirPath, os.ModePerm) // 创建文件保存目录
			if err != nil {
				handleErr(c, 500, err)
				return
			}
			// 5.2 其他文件则直接保存到本地
			err = c.SaveUploadedFile(fileHeader, fileDirPath+"/"+fileHeader.Filename)
			if err != nil {
				handleErr(c, 500, err)
				return
			}
			accessPaths = append(accessPaths, fileAccessDirPath+"/"+fileHeader.Filename) // 添加访问路劲
		}
	}
	c.JSON(200, Result{200, "ok", accessPaths})
}

// QueryListHandle
// @Summary     获取服务器上普通文件或图片的访问列表
// @Description  获取服务器上普通文件或图片的访问列表
// @Tags	filemanage
// @Param        fileType    query     string  true  "文件类型：file or image"
// @Success      200  {object}   handler.Result{data=[]string}
// @Router       /fileList [get]
func QueryListHandle(c *gin.Context) {
	// 1.获取参数
	fileType := c.Query("fileType")
	if fileType == "image" {
		imageDirPath := config.GetImageConfig().ImageDirPath
		imageAccessPath := config.GetImageConfig().ImageAccessDirPath
		imagePaths, err := getAllFiles(imageDirPath)
		if err != nil {
			handleErr(c, 500, err)
			return
		}
		for index, val := range imagePaths {
			imagePaths[index] = imageAccessPath + strings.TrimPrefix(val, imageDirPath)
		}
		c.JSON(200, &Result{200, "ok", imagePaths})
	} else if fileType == "file" {
		fileDirPath := config.GetFileConfig().FileDirPath
		fileAccessPath := config.GetFileConfig().FileAccessDirPath
		filePaths, err := getAllFiles(fileDirPath)
		if err != nil {
			handleErr(c, 500, err)
			return
		}
		for index, val := range filePaths {
			filePaths[index] = fileAccessPath + strings.TrimPrefix(val, fileDirPath)
		}
		c.JSON(200, &Result{200, "ok", filePaths})
	} else {
		c.JSON(200, &Result{400, "参数fileType为file或者image", nil})
	}
}

// 只获取文件而忽略目录
func getAllFiles(dirPath string) ([]string, error) {
	var files []string
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, fileInfo := range fileInfos {
		realPath := dirPath + "/" + fileInfo.Name()
		if fileInfo.IsDir() {
			nextFiles, err := getAllFiles(realPath)
			if err != nil {
				return nil, err
			}
			files = append(files, nextFiles...)
		} else {
			files = append(files, realPath)
		}
	}
	return files, nil
}

// DeleteFileListHandler
// @Summary    删除文件列表
// @Description  服务器上普通文件或图片删除
// @Tags	filemanage
// @Param        fileUrls    body     string  true  "文件的访问url1,文件的访问url2..."
// @Success      200  {object}   handler.Result
// @Router       /deleteFileList [post]
func DeleteFileListHandler(c *gin.Context) {
	fileUrlStr, b := c.GetPostForm("fileUrls")
	if !b {
		c.JSON(200, &Result{Code: 400, Msg: "fileUrls参数不存在"})
		c.Abort()
		return
	}
	if fileUrlStr != "" {
		fileUrls := strings.Split(fileUrlStr, ",")
		if len(fileUrls) > 0 {
			var count = 0 //删除文件数量
			var msg = ""
			for _, fileUrl := range fileUrls {
				// 将url转为文件路径
				var filePath = ""
				if strings.HasPrefix(fileUrl, config.GetFileConfig().FileAccessDirPath) {
					filePath = config.GetFileConfig().FileDirPath + "/" +
						strings.TrimPrefix(fileUrl, config.GetFileConfig().FileAccessDirPath)
				} else if strings.HasPrefix(fileUrl, config.GetImageConfig().ImageAccessDirPath) {
					filePath = config.GetImageConfig().ImageDirPath + "/" +
						strings.TrimPrefix(fileUrl, config.GetImageConfig().ImageAccessDirPath)
				}
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					msg += fmt.Sprintf("\t%s对应的filepath：%s不存在", fileUrl, filePath)
					continue
				}
				if os.Remove(filePath) != nil {
					msg += fmt.Sprintf("\t%s对应的filepath：%s删除失败", fileUrl, filePath)
					continue
				}
				count++
			}
			msg += fmt.Sprintf("\t删除成功：%d个", count)
			c.JSON(200, &Result{Code: 200, Msg: msg})
		}
	} else {
		c.JSON(200, &Result{Code: 400, Msg: "fileUrls长度为0"})
	}
}

// 处理异常，统一返回
func handleErr(c *gin.Context, code int, err error) {
	c.JSON(code, &Result{code, err.Error(), nil})
}
