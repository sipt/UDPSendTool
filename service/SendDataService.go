package service

import (
	"fmt"
	"sipt/UDPSendTool/io"
	"sipt/UDPSendTool/netio"
	"sipt/UDPSendTool/static"
	"strconv"
	"time"
)

const (
	ZERO_STRING = "00000000000000000000"
)

type SendDataService struct {
	FileName     string
	SendFileInfo *io.MyFileInfo
	Count        int64
}

func (this *SendDataService) SetFileName(fileName string) {
	fmt.Println("发送文件前初始化。。。")
	this.FileName = fileName
	this.SendFileInfo = io.GetFileInfo(fileName)
	this.Count = 0
	fmt.Println("初始化完成！", this.SendFileInfo.PackageCount)
}

func (this *SendDataService) Service(netIo *netio.NetIo, fileOp *io.FileOp) {
	//启动发送协程
	go netIo.Send()
	//启动监听协程
	go netIo.Listener()
	//发送传输协议
	this.sendAgreement(netIo)

	this.Count = 0
	go fileOp.ReadFile(this.FileName)
	for {
		temp := <-fileOp.ReadData
		if temp == nil {
			netIo.SendData <- nil
			break
		}
		var bytes []byte
		bytes = []byte(FormatString(this.Count, this.SendFileInfo.PackageNumberLength))
		bytes = append(bytes, temp...)
		fmt.Println(this.Count)
		netIo.SendData <- bytes
		this.Count += 1
	}
}

func (this *SendDataService) sendAgreement(netIo *netio.NetIo) {

	fmt.Println("发送文件协议信息！")
	agreement := this.SendFileInfo.ToString()
	netIo.SendData <- []byte(agreement)
	myTimer := time.NewTimer(time.Second * 2)
	fmt.Println("等待确认信息...")
	for {
		select {
		case <-myTimer.C:
			fmt.Println("等待超时，重新发送")
			netIo.SendData <- []byte(agreement)
			myTimer.Reset(time.Second * 2)
			fmt.Println("等待确认信息...")
		case result := <-netIo.ReadData:
			if string(result) == static.STATUS_SUCCESS {
				fmt.Println("客户端接收成功！")
			} else {
				panic("客户端接收失败！发送取消。。。")
			}
			return
		}
	}
}

func FormatString(count, length int64) string {
	temp := ZERO_STRING + strconv.FormatInt(count, 10)
	return temp[int64(len(temp))-length:]
}
