package Accepter

import "net"

const (
	port      = ":10010"
	blockSize = 4096
)

type DataRead struct {
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
