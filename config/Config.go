package config

import (
	"gopkg.in/ini.v1"
	"log"
)

type Config struct {
	Port          int      `ini:"port"`          // 端口号
	DirRootPath   string   `ini:"dirRootPath"`   // 文件存放根目录
	DirAccessPath string   `ini:"dirAccessPath"` // nginx配置的此根目录的访问路径
	Suffix        []string `ini:"suffix"`        // 图片后缀
	Width         int      `ini:"fixWidth"`      // 压缩后的宽度
}

func LoadConfig() *Config {
	configFile, err := ini.Load("./config/my.ini")
	if err != nil {
		log.Fatalln(err)
	}
	config := &Config{}
	err = configFile.Section("server").MapTo(config)
	if err != nil {
		log.Fatalln(err)
	}
	err = configFile.Section("file").MapTo(config)
	if err != nil {
		log.Fatalln(err)
	}
	err = configFile.Section("image").MapTo(config)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("从my.ini中加载如下配置：%+v\n", *config)
	return config
}
