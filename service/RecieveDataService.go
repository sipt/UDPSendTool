package service

import (
	"fmt"
	"sipt/UDPSendTool/io"
	"sipt/UDPSendTool/netio"
	"sipt/UDPSendTool/static"
	"strconv"
)

type RecivieveDataService struct {
	FileName        string
	RecieveFileInfo *io.MyFileInfo
	RecieveDataList map[int64][]byte
	count           int64
}

func (this *RecivieveDataService) SetFileName(fileName string) {
	this.FileName = fileName
}

func (this *RecivieveDataService) Service(netIo *netio.NetIo, fileOp *io.FileOp) {
	this.count = 0
	firstRead := true
	go netIo.Send()
	go netIo.Listener()
	for {
		bytes := <-netIo.ReadData
		fmt.Println(string(bytes), ".....", firstRead)
		if firstRead && string(bytes[:5]) == "title" {
			this.dealWithData(netIo, bytes)
		} else {
			//存入数据包编号
			id, err := strconv.ParseInt(string(bytes[:this.RecieveFileInfo.PackageNumberLength]), 10, 64)
			if static.DealWithError("包编号转成数字出错", err) {
				//包编号解析出错，丢掉该包
				continue
			}

			//将读入的数据存入map
			this.saveData(id, bytes[this.RecieveFileInfo.PackageNumberLength:])
			//设置协议头接收完毕开始数据传递标志，减轻前if判断压力
			firstRead = false
			//发送确认数据接收包到服务端
			this.sendPackageConfirm(netIo, id)
			//查看数据是否全部读入完毕
			if bytes == nil {
				break
			}
		}
	}
	go fileOp.WriteFile(this.FileName)
	for i := int64(0); i < this.RecieveFileInfo.PackageCount; i++ {
		//按顺序传入文件
		fileOp.WriteData <- this.RecieveDataList[i]
	}
	//传入文件接收完毕标记nil
	fileOp.WriteData <- nil
}

//处理协议接收
func (this *RecivieveDataService) dealWithData(netIo *netio.NetIo, bytes []byte) {
	data := string(bytes)
	this.RecieveFileInfo = io.GetFileInfoByArr(data)
	if this.RecieveFileInfo != nil {
		//发送接收成功标记
		netIo.SendData <- []byte(static.STATUS_SUCCESS)
	} else {
		//发送接收失败标记
		netIo.SendData <- []byte(static.STATUS_FAIL)
	}
}

func (this *RecivieveDataService) getData(bytes []byte) {

}

func (this *RecivieveDataService) saveData(id int64, bytes []byte) {
	this.RecieveDataList[id] = bytes
}

//发送数据包接收完毕信息
func (this *RecivieveDataService) sendPackageConfirm(netIo *netio.NetIo, id int64) {
	netIo.SendData <- []byte(strconv.FormatInt(id, 10))
}
