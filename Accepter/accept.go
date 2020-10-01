package Accepter

import (
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"strconv"
)

var dataRead DataRead
var dataPack DataPack
var host Host

func (h Host) Connect() *net.TCPConn {
	host, _ := net.ResolveTCPAddr("tcp4", h.IP+h.Port)
	fmt.Println("监听：", host.IP, host.Port)
	listener, err := net.ListenTCP("tcp", host)
	if err != nil {
		color.Red("监听失败", err)
	}
	conn, err := listener.AcceptTCP()
	if err != nil {
		color.Red("接收客户端失败", err)
	}
	return conn
}

func Accept() bool {
	host.IP = "0.0.0.0"
	host.Port = port
	conn := host.Connect()
	defer conn.Close()
	AcceptInfo(conn)

	fileLen := int(dataRead.FileLen)
	for i := 0; i < fileLen; i++ {
		dataPack.UnPack(conn)
		dataPack.Receive(conn)
	}
	return true
}

func AcceptInfo(conn *net.TCPConn) bool {
	dataRead.SizeTotal = AcceptSize(conn)
	if dataRead.SizeTotal == 0 {
		color.Red("接收文件总大小有误")
		return false
	}
	dataRead.FileLen = AcceptSize(conn)
	if dataRead.FileLen == 0 {
		color.Red("接收文件数量有误")
		return false
	}
	dataRead.Folder = AcceptPath(conn)
	if len(dataRead.Folder) != 0 {
		for n, tmp := 1, dataRead.Folder; IsExit(dataRead.Folder); {
			dataRead.Folder = tmp + "-副本" + strconv.Itoa(n)
			n++
		}
		_ = os.MkdirAll(dataRead.Folder, os.ModePerm)
	}

	return true
}

func (dp *DataPack) UnPack(conn *net.TCPConn) {
	//dp.Connect = host.Connect()
	dp.FilePath = AcceptPath(conn)
	dp.FileSize = int(AcceptSize(conn))
}

func (dp DataPack) Receive(conn *net.TCPConn) bool {
	data := make(chan []byte, blockSize)
	writerResult := make(chan bool)
	receiveResult := make(chan bool)
	counter := make(chan int64)
	go func() {
		writerResult <- FileWriter(dp.FilePath, data)
	}()
	go func() {
		receiveResult <- Receiver(conn, data, true, counter)
	}()

	go DisplayCounterAccept(dataRead.SizeTotal, int64(dp.FileSize), counter)
	if <-writerResult && <-receiveResult {
		return true
	} else {
		color.Red("接收文件失败")
		return false
	}
}

func AcceptPath(conn *net.TCPConn) string {
	tmp := make([]byte, 200)
	n, err := conn.Read(tmp)
	if err != nil {
		color.Red("接收文件路径失败", err)
		tmp = []byte("fail")
		_, _ = conn.Write(tmp)
		return ""
	}
	res := string(tmp[:n])
	tmp = []byte("success")
	_, _ = conn.Write(tmp)
	return res
}

func AcceptSize(conn *net.TCPConn) int64 {
	tmp := make([]byte, 200)
	n, err := conn.Read(tmp)
	if err != nil {
		color.Red("接收数据失败", err)
		tmp = []byte("fail")
		_, _ = conn.Write(tmp)
		return 0
	}
	res, _ := strconv.ParseInt(string(tmp[:n]), 10, 64)
	tmp = []byte("success")
	_, _ = conn.Write(tmp)
	return res
}

var acceptSize int64
func DisplayCounterAccept(totalSize int64, fileSize int64, counter chan int64) {
	var tmpSize = int64(0)
	green := color.New(color.FgGreen)
	for tmp := range counter {
		tmpSize += tmp
		if tmpSize > fileSize {
			tmp = fileSize % blockSize
		}
		acceptSize += tmp
		_, _ = green.Printf("总进度：%d%%\r", int(float64(acceptSize)/float64(totalSize)*100))
	}
	if acceptSize == totalSize {
		color.Yellow("\n所有文件已接收！")
	}
}
