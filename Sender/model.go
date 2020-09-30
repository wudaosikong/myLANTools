package Sender

import "net"

const port string = ":10086"

type DataRead struct {
	PathArry  []string
	SizeArry  []int64
	SizeTotal int64
	FileLen   int64
}
type DataPack struct {
	FilePath string
	FileSize int64
	Connect  *net.TCPConn
}

type Host struct {
	IP   string
	Port string
}
