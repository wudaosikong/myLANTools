package Sender

const (
	port      = ":10086"
	blockSize = 4096
)

type DataRead struct {
	PathArray []string
	SizeArray []int64
	SizeTotal int64
	FileLen   int64
	Folder	string
}
type DataPack struct {
	FilePath string
	FileSize int64
}

type Host struct {
	IP   string
	Port string
}
