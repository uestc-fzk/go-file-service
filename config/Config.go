package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Config struct {
	Port          string `properties:"port"`          // 端口号
	DirRootPath   string `properties:"dirRootPath"`   // 文件存放根目录
	DirAccessPath string `properties:"dirAccessPath"` // nginx配置的此根目录的访问路径
}

var fieldMap = make(map[string]string) // 从属性的标签中获取的别名alias映射指向属性名，即fieldAlias-->fieldName

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 1.反射获取Type和Value
	var config Config
	ty := reflect.TypeOf(config)
	va := reflect.ValueOf(&config) // 这里必须是指针才能修改其属性

	// 2.构建fieldMap
	for i := 0; i < ty.NumField(); i++ {
		field := ty.Field(i)
		var tag = (string)(field.Tag)
		splits := strings.Split(tag, ",")
		for _, str := range splits {
			// 找到properties前缀并取值
			if strings.HasPrefix(str, "properties") {
				// 取出引号中的值，即fieldAlias
				fieldAlias := str
				fieldAlias = fieldAlias[strings.IndexRune(str, '"')+1:]
				fieldAlias = fieldAlias[:strings.IndexRune(fieldAlias, '"')]
				// 放入map中
				fieldMap[fieldAlias] = field.Name
			}
		}
	}
	fmt.Printf("%+v\n", fieldMap)

	// 3.读取配置文件
	file, err := os.Open("./config/config.properties")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	// 4.循环读取每一行，并反射修改属性值
	for {
		lineStr, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			if lineStr != "" {
				changeFieldVal(lineStr, &va)
			}
			break
		}
		changeFieldVal(lineStr, &va)
	}
	// 5.打印结果并返回
	fmt.Printf("%+v \n", config)
	return &config
}

// 修改属性值
func changeFieldVal(lineStr string, vaPtr *reflect.Value) {
	va := *vaPtr
	str := lineStr[:strings.IndexRune(lineStr, ' ')]
	fieldAlias := str[:strings.IndexRune(str, '=')]
	fieldVal := str[strings.IndexRune(str, '=')+1:]
	fieldName, ok := fieldMap[fieldAlias]
	if ok {
		field := va.Elem().FieldByName(fieldName)
		field.SetString(fieldVal) // 反射设置值
	}
}
