package Sender

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

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
