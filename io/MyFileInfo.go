package io

import (
	"fmt"
	"os"
)

type MyFileInfo struct {
	FileSize            int64
	FileName            string
	PackageCount        int64
	PackageNumberLength int64
}

func (this *MyFileInfo) GetFileInfo(fileName string) {
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("读取文件信息时出错误！")
		panic(err)
	}
	this.FileSize = fileInfo.Size()
	this.FileName = fileName
	this.PackageCount
}
