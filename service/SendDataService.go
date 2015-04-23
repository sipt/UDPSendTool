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
	SendDataList map[int64][]byte
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
	go this.removeFromMapById(netIo)
	for {
		//从通道里读出数据
		temp := <-fileOp.ReadData

		if temp == nil {
			//第一遍数据发送完毕，开发把客户端示收到数据重发
			for {
				//查看是否存在重发数据
				if len(this.SendDataList) > 0 {
					//开始重发操作
					for _, value := range this.SendDataList {
						netIo.SendData <- value
					}
				} else {
					break
				}
			}
			//传给客户文件传输结束标记
			netIo.SendData <- nil
			break
		}
		var bytes []byte
		//格式化计数到指定位数
		bytes = []byte(FormatString(this.Count, this.SendFileInfo.PackageNumberLength))
		//在后面补上数据
		bytes = append(bytes, temp...)
		//将数据保存入map队列
		this.SendDataList[this.Count] = bytes
		//将格式化的数据发送出去
		netIo.SendData <- bytes
		this.Count += 1
	}
}

//用于监听客户发来的确认包，并从map中移除
func (this *SendDataService) removeFromMapById(netIo *netio.NetIo) {
	for {
		idBytes := <-netIo.ReadData
		id, err := strconv.ParseInt(string(idBytes), 10, 64)
		if static.DealWithError("接收到客户端确认包，包编号转成数字出错", err) {
			//丢弃该包
			continue
		}
		//删除客户已经接收到的数据
		delete(this.SendDataList, id)
	}

}

//发送文件的协议
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
