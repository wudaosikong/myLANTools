package Accepter

import (
	"io"
	"net"

	"github.com/fatih/color"
)

func Receiver(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(data)
	defer close(counter)

	var tmpSize = 0

	for {
		tmp := make([]byte, blockSize)
		n, err := conn.Read(tmp)
		tmpSize += n
		if err != nil && err != io.EOF {
			color.Red("接收失败", err)
			return false
		} else if err == io.EOF {
			return true
		}
		if tmpSize > dataPack.FileSize {
			n = dataPack.FileSize % blockSize

			data <- tmp[:n]

			if isDisplay {
				counter <- int64(n)
			}
			return true
		}

		data <- tmp[:n]

		if isDisplay {
			counter <- int64(n)
		}
	}
}
