package service

import (
	"demo/client/io"
	"demo/client/netio"
)

type IDataService interface {
	SetFileName(string)
	Service(*netio.NetIo, *io.FileOp)
}
