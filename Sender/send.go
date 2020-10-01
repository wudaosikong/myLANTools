package Sender

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fatih/color"
)

func (h Host) Connect() *net.TCPConn {
	host, _ := net.ResolveTCPAddr("tcp4", h.IP+h.Port)
	conn, err := net.DialTCP("tcp", nil, host)
	if err != nil {
		color.Red("连接对方主机失败", err)
	}
	return conn
}

var dataRead DataRead
var dataPack DataPack
var host Host
var start time.Time
func Send(dir string, ip string) bool {
	host.IP = ip
	host.Port = port

	dataRead.Read(dir, ip)
	conn := host.Connect()
	start=time.Now()
	defer conn.Close()
	SendInfo(conn)

	for i, _ := range dataRead.PathArray {
		dataPack.Pack(i)
		dataPack.Send(conn)
	}
	return true
}

func SendInfo(conn *net.TCPConn) bool {
	if !SendSize(dataRead.SizeTotal, conn) {
		return false
	}
	if !SendSize(dataRead.FileLen, conn) {
		return false
	}
	if !SendPath(dataRead.Folder, conn) {
		return false
	}
	return true
}

func (dr *DataRead) Read(dir string, ip string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			tmp := filepath.ToSlash(path)
			dr.PathArray = append(dr.PathArray, tmp)
			dr.SizeArray = append(dr.SizeArray, info.Size())
			dr.SizeTotal += info.Size()
		}
		return nil
	})
	dr.FileLen = int64(len(dataRead.PathArray))
	if info, _ := os.Stat(dir); info.IsDir() {
		dr.Folder = dir
	}
}

func (dp *DataPack) Pack(i int) {
	//dp.Connect = host.Connect()
	dp.FilePath = dataRead.PathArray[i]
	dp.FileSize = dataRead.SizeArray[i]
}

func (dp DataPack) Send(conn *net.TCPConn) bool {
	if !SendPath(dp.FilePath, conn) {
		return false
	}
	if !SendSize(dp.FileSize, conn) {
		return false
	}

	readerResult := make(chan bool)
	senderResult := make(chan bool)
	counter := make(chan int64)
	data := make(chan []byte, blockSize)
	go func() {
		readerResult <- FileReader(dp.FilePath, data)
	}()
	go func() {
		senderResult <- Sender(conn, data, true, counter)
	}()

	go DisplayCounterSend(dataRead.SizeTotal, dp.FileSize, counter)

	if <-readerResult && <-senderResult {
		return true
	} else {
		color.Red("发送失败")
		return false
	}
}

func SendPath(path string, client *net.TCPConn) bool {
	tmpName := []byte(path)
	_, err := client.Write(tmpName)
	if err != nil {
		color.Red("发送文件路径失败", err)
		return false
	}
	tmp := make([]byte, 7)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件(夹)名失败")
		return false
	}
	return true
}

func SendSize(size int64, client *net.TCPConn) bool {
	tmpSize := make([]byte, 200)
	tmpSize = []byte(strconv.FormatInt(size, 10))
	_, err := client.Write(tmpSize)
	if err != nil {
		color.Red("发送文件大小失败", err)
		return false
	}
	tmp := make([]byte, 7)
	n, _ := client.Read(tmp)
	if string(tmp[:n]) != "success" {
		color.Red("对方接收文件大小失败")
		return false
	}
	return true
}

var sendSize int64
func DisplayCounterSend(totalSize int64, fileSize int64, counter chan int64) {
	var tmpSize = int64(0)
	green := color.New(color.FgGreen)
	for tmp := range counter {
		tmpSize += tmp
		if tmpSize > fileSize {
			tmp = fileSize % blockSize
		}
		sendSize += tmp
		_, _ = green.Printf("总进度：%d%%\r", int(float64(sendSize)/float64(totalSize)*100))
	}
	if sendSize == totalSize {
		cost := time.Since(start)
		color.Yellow("\n所有文件已发送！\n")
		fmt.Printf("耗时：%d s\n速度：%.2f Mb/s",int(cost.Seconds()),float64(totalSize)/cost.Seconds()/1024/1024)
	}
}
