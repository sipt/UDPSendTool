package netio

import (
	"fmt"
	"net"
)

type NetIo struct {
	SendData    chan []byte
	ReadData    chan []byte
	SendIp      [4]byte
	ListenIp    [4]byte
	SendToPort  int
	ListenPort  int
	IsWriting   bool
	ReadEndChan chan string
	// fileOp    *io.FileOp
}

func (this *NetIo) Send() {
	socket, err := net.DialUDP("udp4", nil, &net.UDPAddr{
		IP:   net.IPv4(this.SendIp[0], this.SendIp[1], this.SendIp[2], this.SendIp[3]),
		Port: this.SendToPort,
	})
	if err != nil {
		fmt.Println("连接失败！", err)
		return
	}
	defer socket.Close()
	for {
		//发送数据
		senddata := <-this.SendData
		_, err = socket.Write(senddata)
		if err != nil {
			fmt.Println("发送数据失败!", err)
			return
		}
		if senddata == nil {
			return
		}
	}
}

func (this *NetIo) Listener() {
	this.ReadEndChan = make(chan string)
	this.IsWriting = false
	//创建监听
	socket, err := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: this.ListenPort,
	})
	if err != nil {
		fmt.Println("监听失败", err)
		return
	}
	defer socket.Close()
	for {
		//读取数据
		data := make([]byte, 2048)
		length, _, err := socket.ReadFromUDP(data)
		if err != nil {
			fmt.Println("读取数据失败！", err)
			continue
		}
		fmt.Println("--------读入数据长度:", length)
		if length == 0 {
			this.ReadEndChan <- "END"
		} else {
			this.ReadData <- data[:length]
		}
	}
}
