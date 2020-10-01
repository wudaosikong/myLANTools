package Accepter

const (
	port      = ":10086"
	blockSize = 4096
)

type DataRead struct {
	SizeTotal int64
	FileLen   int64
	Folder	string
}
type DataPack struct {
	FilePath string
	FileSize int
}

type Host struct {
	IP   string
	Port string
}
