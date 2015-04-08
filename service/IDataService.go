package service

import (
	"sipt/UDPSendTool/io"
	"sipt/UDPSendTool/netio"
)

type IDataService interface {
	SetFileName(string)
	Service(*netio.NetIo, *io.FileOp)
}
