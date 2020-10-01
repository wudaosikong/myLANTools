package Accepter

import (
	"fmt"
	"os"
	"strings"
)

// 判断文件是否存在
func IsExit(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func FileWriter(filename string, data chan []byte) bool {
	var tmpDir string
	if strings.Index(filename, "/") != -1 &&len(dataRead.Folder)!=0{
		filename = dataRead.Folder + filename[strings.Index(filename, "/"):]
	}
	if strings.LastIndex(filename, "/") != -1{
		tmpDir=filename[:strings.LastIndex(filename, "/")]
		if !IsExit(tmpDir){
			os.MkdirAll(tmpDir,os.ModePerm)
		}
	}
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("文件创建失败", err)
		return false
	}
	defer file.Close()
	var tmpSize = 0
	for bytes := range data {
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println("文件写入错误", err)
			return false
		}
		if tmpSize == dataPack.FileSize {
			return true
		}
	}
	return true
}
