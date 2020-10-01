package Sender

import (
	"net"

	"github.com/fatih/color"
)

func Sender(conn *net.TCPConn, data chan []byte, isDisplay bool, counter chan int64) bool {
	defer close(counter)

	for tmp := range data {
		_, err := conn.Write(tmp)
		if err != nil {
			color.Red("发送失败", err)
			return false
		}
		if isDisplay {
			counter <- int64(len(tmp))
		}
	}

	return true
}
