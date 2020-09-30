package Sender

import (
	"fmt"
	"myLANTools/Share"
	"net"
	"os"
	"path/filepath"
	"strconv"

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

func Send(dir string, ip string) bool {
	host.IP = ip
	host.Port = port

	dataRead.Read(dir, ip)
	SendInfo()

	for i, _ := range dataRead.PathArry {
		dataPack.Pack(i)
		dataPack.Send()
	}
	return true
}

func SendInfo() bool {
	conn := host.Connect()
	defer conn.Close()
	if !SendSize(dataRead.SizeTotal, conn) {
		return false
	}
	if !SendSize(dataRead.FileLen, conn) {
		return false
	}
	return true
}

func (dr *DataRead) Read(dir string, ip string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			dr.PathArry = append(dr.PathArry, path)
			dr.SizeArry = append(dr.SizeArry, info.Size())
			dr.SizeTotal += dr.SizeTotal
		}
		return nil
	})
	dr.FileLen = int64(len(dataRead.PathArry))
}

func (dp *DataPack) Pack(i int) {
	dp.Connect = host.Connect()
	dp.FilePath = dataRead.PathArry[i]
	dp.FileSize = dataRead.SizeArry[i]
}

func (dp DataPack) Send() bool {
	if !SendPath(dp.FilePath, dp.Connect) {
		return false
	}
	if !SendSize(dp.FileSize, dp.Connect) {
		return false
	}

	readerResult := make(chan bool)
	senderResult := make(chan bool)
	counter := make(chan int64)
	data := make(chan []byte, 1024)
	go func() {
		readerResult <- Share.FileReader(dp.FilePath, data)
	}()
	go func() {
		senderResult <- Share.Sender(dp.Connect, data, true, counter)
	}()

	go DisplayCounterSend(dataRead.SizeTotal, counter)

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
		color.Red("发送文件(夹)名失败", err)
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

func DisplayCounterSend(size int64, counter chan int64) {
	var sendSize int64
	green := color.New(color.FgGreen)
	for tmp := range counter {
		sendSize += tmp
		_, _ = green.Printf("总进度：%d%%\r", int(float64(sendSize)/float64(size)*100))
	}
	fmt.Println("")
}

func isFile(dir string) bool {
	f, _ := os.Stat(dir)
	if f.IsDir() {
		return false
	}
	return true
}
