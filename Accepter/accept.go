package Accepter

import (
	"fmt"
	"net"
	"strconv"

	"github.com/fatih/color"
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

	AcceptInfo()

	fileLen := int(dataRead.FileLen)
	for i := 0; i < fileLen; i++ {

	}
	return true
}

func AcceptInfo() bool {
	conn := host.Connect()
	defer conn.Close()
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
	return true
}

func FileReceive(filename string, conn *net.TCPConn, size int64) bool {
	data := make(chan []byte, blockSize)
	writerResult := make(chan bool)
	receiveResult := make(chan bool)
	counter := make(chan int64)
	go func() {
		writerResult <- FileWriter(filename, data)
	}()
	go func() {
		receiveResult <- Receiver(conn, data, true, counter)
	}()

	go DisplayCounterAccept(size, counter)
	if <-writerResult && <-receiveResult {
		return true
	} else {
		color.Red("接收文件失败")
		return false
	}
}

func (dp *DataPack) UnPack(i int) {
	dp.Connect = host.Connect()
	dp.FilePath = dataRead.PathArry[i]
	dp.FileSize = dataRead.SizeArry[i]
}

func AcceptPath(conn *net.TCPConn) string {
	tmp := make([]byte, 200)
	n, err := conn.Read(tmp)
	if err != nil {
		color.Red("接收文件信息&文件名失败", err)
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
