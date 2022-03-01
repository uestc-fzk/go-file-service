package handler

import (
	"GoFileService/config"
	"fmt"
	"github.com/nfnt/resize"
	"gopkg.in/ini.v1"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// CompressImage 将源图片压缩到目标图片
func CompressImage(imgSrcPath string, imgDestPath string) error {
	// 1.打开源图片
	srcFile, err := os.Open(imgSrcPath)
	if err != nil {
		return err
	}
	defer func() {
		_ = srcFile.Close()
		_ = os.Remove(imgSrcPath) // 移除原文件
	}()
	// 2.decode源图片
	srcImg, format, err := image.Decode(srcFile)
	if err != nil {
		return err
	}

	fmt.Printf("协程id: %d\t图片格式: %s\t 开始压缩===> %s \n", GoroutineId(), format, imgDestPath)
	// 3.缩放操作
	destImg := resize.Resize(uint(config.GetImageConfig().FixWidth), 0, srcImg, resize.Lanczos3)
	// 4.新建目的图片
	destFile, err := os.Create(imgDestPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// 5.缩放源图片至目的图片: 这里只压缩jpeg和png格式的图片
	if format == "png" {
		err = png.Encode(destFile, destImg)
	} else if format == "jpeg" { // jpg会被golang解析为jpeg格式
		// 100%质量的保存
		err = jpeg.Encode(destFile, destImg, &jpeg.Options{Quality: 100})
	} else {
		log.Printf("只能缩放png或者jpeg/jpg格式图片，此图片格式为%s，其源图片%s将被自动删除\n", format, srcFile.Name())
	}
	if err != nil {
		return err
	}
	return nil
}

// IsImage 判断是否是图片
func IsImage(fileName string) bool {
	splits := strings.SplitAfter(fileName, ".")
	// 说明没有格式，肯定不是图片
	if len(splits) <= 1 {
		return false
	}
	// 取出格式
	var format = strings.ToLower(splits[len(splits)-1])
	_, ok := ImageSuffixMap[format]
	return ok
}

// CanComPress 判断是否能压缩，即判断是否为png或者jpeg/jpg格式
func CanComPress(fileName string) bool {
	splits := strings.SplitAfter(fileName, ".")
	// 说明没有格式，肯定不是图片
	if len(splits) <= 1 {
		return false
	}
	// 取出格式
	var format = strings.ToLower(splits[len(splits)-1])
	return format == "jpg" || format == "jpeg" || format == "png"
}

var ImageSuffixMap = make(map[string]int, 16)

func init() {
	configFile, err := ini.Load("./config/my.ini")
	if err != nil {
		log.Fatalln(err)
	}
	for _, suffix := range configFile.Section("image").Key("suffix").Strings(",") {
		ImageSuffixMap[suffix] = 1
	}
	fmt.Printf("图片压缩处理器init...%+v\n", ImageSuffixMap)
}

func GoroutineId() int {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic recover:panic info:%v\n", err)
		}
	}()

	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	}
	return id
}
