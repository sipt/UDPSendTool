package service

import (
	"demo/client/io"
	"demo/client/netio"
)

type RecivieveDataService struct {
	FileName string
}

func (this *RecivieveDataService) SetFileName(fileName string) {
	this.FileName = fileName
}

func (this *RecivieveDataService) Service(netIo *netio.NetIo, fileOp *io.FileOp) {
	go netIo.Listener()
	go fileOp.WriteFile(this.FileName)
	for {
		bytes := <-netIo.ReadData
		fileOp.WriteData <- bytes
		if bytes == nil {
			break
		}
	}
}
