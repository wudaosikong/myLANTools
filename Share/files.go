package Share

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const blockSize = 1024

// 判断文件是否存在
func IsExit(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	return false
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func FileReader(filename string, data chan []byte) bool {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开文件出错", err)
		return false
	}
	defer file.Close()
	defer close(data)

	reader := bufio.NewReader(file)
	for {
		tmp := make([]byte, blockSize)
		n, err := reader.Read(tmp)
		if err != nil && err != io.EOF {
			fmt.Println("文件读取错误", err)
			return false
		}
		if n == 0 {
			return true
		}
		data <- tmp
	}
}

func FileWriter(filename string, data chan []byte) bool {
	var tmpDir string
	if strings.LastIndex(filename, "/") != -1 {
		tmpDir = filename[:strings.LastIndex(filename, "/")]
	}
	if !IsExit(tmpDir) {
		_ = os.MkdirAll(tmpDir, os.ModePerm)
	}
	for n, tmp := 1, filename; IsExit(filename); {
		filename = tmp[:strings.LastIndex(tmp, ".")] + "-副本" + strconv.Itoa(n) + tmp[strings.LastIndex(tmp, "."):]
		n++
	}
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("文件创建失败", err)
		return false
	}
	defer file.Close()
	for bytes := range data {
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println("文件写入错误", err)
			return false
		}
	}
	return true
}
