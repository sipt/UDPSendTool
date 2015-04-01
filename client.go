package main

import (
	"demo/client/io"
	"demo/client/netio"
	"fmt"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// fmt.Print("请输入IP地址：")
	// fmt.Scanf("%d.%d.%d.%d", &netIo.SendIp[0], &netIo.SendIp[1], &netIo.SendIp[2], &netIo.SendIp[3])
	// fmt.Print("请输入端口号：")
	// fmt.Scanf("\n%d", &netIo.Port)
	// fmt.Println(netIo.SendIp, ":", netIo.Port)
	netIo := &netio.NetIo{}
	fileOp := &io.FileOp{}
	var order int
	fmt.Print("请输入指令：")
	fmt.Scanf("%d", &order)
	fmt.Println(order)
	if order == 0 {
		netIo.SendIp[0] = 127
		netIo.SendIp[1] = 0
		netIo.SendIp[2] = 0
		netIo.SendIp[3] = 1
		netIo.Port = 10000
		netIo.SendData = make(chan []byte, 5)
		fileOp.ReadData = make(chan []byte, 5)
		fmt.Println(netIo.SendIp)
		fmt.Println(netIo.Port)
		go netIo.Send()
		go fileOp.ReadFile("D:\\GOPATH\\src\\demo\\1.png")
		for {
			bytes := <-fileOp.ReadData
			netIo.SendData <- bytes
			if bytes == nil {
				break
			}
		}
	} else {

		netIo.Port = 10000
		netIo.ReadData = make(chan []byte, 100)
		fileOp.WriteData = make(chan []byte, 5)
		go netIo.Listener()
		go fileOp.WriteFile("D:\\2.png")
		for {
			bytes := <-netIo.ReadData
			fileOp.WriteData <- bytes
			if bytes == nil {
				break
			}
		}
	}
}
