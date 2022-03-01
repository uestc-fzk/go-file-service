package handler

import (
	"github.com/nfnt/resize"
	"gopkg.in/ini.v1"
	"image"
	"image/jpeg"
	"log"
	"os"
	"strings"
)

func HandleImageFile(dirPath string, fileName string) error {
	file, err := os.Open(dirPath + "/" + fileName)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
		_ = os.Remove(dirPath + "/" + fileName) // 移除原文件
	}()

	fileImg, _, err := image.Decode(file)
	if err != nil {
		return err
	}
	err = resizeToFixWidth(&fileImg, fixWidth, dirPath+"/_"+fileName)
	if err != nil {
		return err
	}
	return nil
}

// resizeToFixWidth 需要确保父目录已经创建
func resizeToFixWidth(srcImg *image.Image, fixWidth int, filePath string) error {
	log.Printf("图片压缩===> %s \n", filePath)
	newImg := resize.Resize(uint(fixWidth), 0, *srcImg, resize.Lanczos3)
	imgFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer imgFile.Close()
	// 100%质量的保存
	err = jpeg.Encode(imgFile, newImg, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	return nil
}

// IsImage 判断是否是图片
func IsImage(fileName string) bool {
	strings.Split(fileName, ".")
	splits := strings.SplitAfter(fileName, ".")
	// 说明没有格式，肯定不是图片
	if len(splits) <= 1 {
		return false
	}
	// 取出格式
	var format = splits[len(splits)-1]
	_, ok := ImageSuffixMap[format]
	log.Println(format, ok)
	return ok
}

var ImageSuffixMap = make(map[string]int, 16)
var fixWidth = 0

func init() {
	configFile, err := ini.Load("./config/my.ini")
	if err != nil {
		log.Fatalln(err)
	}
	for _, suffix := range configFile.Section("image").Key("suffix").Strings(",") {
		ImageSuffixMap[suffix] = 1
	}
	fixWidth, err = configFile.Section("image").Key("fixWidth").Int()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("图片压缩处理器init...%+v\n压缩图片固定宽度: %d px\n", ImageSuffixMap, fixWidth)
}
