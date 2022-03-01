package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	all := getAllFiles("D:/test")
	for _, v := range all {
		fmt.Println(v)
	}
}

// 只获取文件而忽略目录
func getAllFiles(dirPath string) []string {
	var files []string
	fileInfos, _ := ioutil.ReadDir(dirPath)
	for _, fileInfo := range fileInfos {
		realPath := dirPath + "/" + fileInfo.Name()
		if fileInfo.IsDir() {
			files = append(files, getAllFiles(realPath)...)
		} else {
			files = append(files, realPath)
		}
	}
	return files
}
