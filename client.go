package main

import (
	"fmt"
	"runtime"
	"sipt/UDPSendTool/io"
	"sipt/UDPSendTool/netio"
	"sipt/UDPSendTool/service"
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

	var ids service.IDataService
	netIo.SendIp[0] = 127
	netIo.SendIp[1] = 0
	netIo.SendIp[2] = 0
	netIo.SendIp[3] = 1

	netIo.SendData = make(chan []byte, 5)
	fileOp.ReadData = make(chan []byte, 5)
	netIo.ReadData = make(chan []byte, 100)
	fileOp.WriteData = make(chan []byte, 5)

	if order == 0 {
		netIo.SendToPort = 10001
		netIo.ListenPort = 10000
		ids = service.GetDataService(service.SEND_DATA)
		ids.SetFileName("D:\\GOPATH\\src\\demo\\1.png")
	} else {
		netIo.SendToPort = 10000
		netIo.ListenPort = 10001
		ids = service.GetDataService(service.RECIEVE_DATA)
		ids.SetFileName("D:\\2.png")
	}
	ids.Service(netIo, fileOp)
}
