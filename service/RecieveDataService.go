package service

import (
	"fmt"
	"sipt/UDPSendTool/io"
	"sipt/UDPSendTool/netio"
	"sipt/UDPSendTool/static"
)

type RecivieveDataService struct {
	FileName        string
	RecieveFileInfo *io.MyFileInfo
}

func (this *RecivieveDataService) SetFileName(fileName string) {
	this.FileName = fileName
}

func (this *RecivieveDataService) Service(netIo *netio.NetIo, fileOp *io.FileOp) {
	firstRead := true
	go netIo.Send()
	go netIo.Listener()
	for {
		bytes := <-netIo.ReadData
		fmt.Println(string(bytes), ".....", firstRead)
		if firstRead && string(bytes[:5]) == "title" {
			this.dealWithData(netIo, bytes)
		} else {
			if firstRead {
				go fileOp.WriteFile(this.FileName)
				firstRead = false
			}
			fileOp.WriteData <- bytes
			if bytes == nil {
				break
			}
		}
	}
}

func (this *RecivieveDataService) dealWithData(netIo *netio.NetIo, bytes []byte) {
	data := string(bytes)
	this.RecieveFileInfo = io.GetFileInfoByArr(data)
	if this.RecieveFileInfo != nil {
		netIo.SendData <- []byte(static.STATUS_SUCCESS)
	} else {
		netIo.SendData <- []byte(static.STATUS_FAIL)
	}
}

func (this *RecivieveDataService) getData(bytes []byte) {

}
