package service

import (
	"demo/client/io"
	"demo/client/netio"
)

type SendDataService struct {
	FileName string
}

func (this *SendDataService) SetFileName(fileName string) {
	this.FileName = fileName
}

func (this *SendDataService) Service(netIo *netio.NetIo, fileOp *io.FileOp) {
	go netIo.Send()
	go fileOp.ReadFile(this.FileName)
	for {
		bytes := <-fileOp.ReadData
		netIo.SendData <- bytes
		if bytes == nil {
			break
		}
	}
}
