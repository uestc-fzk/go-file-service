package config

import (
	"fmt"
	"gopkg.in/ini.v1"
	"log"
	"runtime"
)

type ServerConfig struct {
	Port int `ini:"port"`
}
type FileConfig struct {
	FileDirPath        string // 根据当前系统选择的存放文件的目录
	FileDirPathOfWin   string `ini:"fileDirPathOfWin"`   // windows文件存放根目录
	FileDirPathOfLinux string `ini:"fileDirPathOfLinux"` // Linux文件存放根目录
	FileAccessDirPath  string `ini:"fileAccessDirPath"`  // nginx配置的此根目录的访问路径
}
type ImageConfig struct {
	Suffix              []string `ini:"suffix"`         // 图片格式
	FixWidth            int      `ini:"fixWidth"`       // 压缩至此宽度px
	ToCompressSize      int      `ini:"toCompressSize"` // 图片达此限制将被压缩，单位byte
	ImageDirPath        string   // 根据当前系统选择的存放图片的目录
	ImageDirPathOfWin   string   `ini:"imageDirPathOfWin"`   // windows放图片的目录
	ImageDirPathOfLinux string   `ini:"imageDirPathOfLinux"` // linux放图片的目录
	ImageAccessDirPath  string   `ini:"imageAccessDirPath"`  // nginx配置的访问存放图片目录的访问路劲
}

var (
	serverConfig *ServerConfig = new(ServerConfig)
	fileConfig   *FileConfig   = new(FileConfig)
	imageConfig  *ImageConfig  = new(ImageConfig)
)

func init() {
	configFile, err := ini.Load("./config/my.ini")
	if err != nil {
		log.Fatalln(err)
	}
	err = configFile.Section("server").MapTo(serverConfig)
	if err != nil {
		log.Fatalln(err)
	}
	err = configFile.Section("file").MapTo(fileConfig)
	if err != nil {
		log.Fatalln(err)
	}
	err = configFile.Section("image").MapTo(imageConfig)
	if err != nil {
		log.Fatalln(err)
	}

	// 根据当前系统环境，选择文件和图片的存放目录
	if runtime.GOOS == "windows" {
		fileConfig.FileDirPath = fileConfig.FileDirPathOfWin
		imageConfig.ImageDirPath = imageConfig.ImageDirPathOfWin
	} else {
		fileConfig.FileDirPath = fileConfig.FileDirPathOfLinux
		imageConfig.ImageDirPath = imageConfig.ImageDirPathOfLinux
	}
	fmt.Printf("当前系统环境为%s, FileDirPATH选择为：%s， ImageDirPATH选择为：%s\n", runtime.GOOS, fileConfig.FileDirPath, imageConfig.ImageDirPath)
	fmt.Printf("读取my.ini配置文件，解析如下:\n"+
		"ServerConfig: %+v\n"+
		"FileConfig: %+v\n"+
		"ImageConfig: %+v\n", *serverConfig, *fileConfig, *imageConfig)
}

func GetServerConfig() *ServerConfig {
	return serverConfig
}
func GetFileConfig() *FileConfig {
	return fileConfig
}
func GetImageConfig() *ImageConfig {
	return imageConfig
}
