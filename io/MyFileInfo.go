package io

//BUNYWAAALXMpVGYS
import (
	"fmt"
	"os"
	"sipt/UDPSendTool/static"
	"strconv"
	"strings"
)

/*
	通过文件全路径名获取文件信息
*/
func GetFileInfo(fileName string) *MyFileInfo {
	fmt.Println("开始获取要发送文件的信息。。。")

	myFileInfo := &MyFileInfo{}
	fileInfo, err := os.Stat(fileName)
	if err != nil {
		fmt.Println("读取文件信息时出错误！")
		panic(err)
	}
	myFileInfo.FileSize = fileInfo.Size()
	myFileInfo.FileName = fileName
	myFileInfo.PackageCount = myFileInfo.FileSize / static.PACKAGE_SIZE
	if myFileInfo.FileSize%static.PACKAGE_SIZE != 0 {
		myFileInfo.PackageCount += 1
	}
	myFileInfo.PackageNumberLength = GetNumberLength(myFileInfo.PackageCount)

	fmt.Println("获取要发送文件的信息成功！")

	return myFileInfo
}

//获取数字位数
func GetNumberLength(number int64) int64 {
	var count int64
	count = 0
	for number > 0 {
		if number%10 != 0 {
			count += 1
		}
		number /= 10
	}
	return count
}

//通过文件信息的字串还原对象
func GetFileInfoByArr(str string) *MyFileInfo {
	myFileInfo := &MyFileInfo{}
	args := strings.Split(str, "{:}")
	var err error
	for i := 1; i <= 4; i++ {
		err = nil
		temp := strings.Split(args[i], "=")
		fmt.Println(temp)
		switch temp[0] {
		case "filesize":
			myFileInfo.FileSize, err = strconv.ParseInt(temp[1], 10, 64)
		case "filename":
			myFileInfo.FileName = temp[1]
		case "packagecount":
			myFileInfo.PackageCount, err = strconv.ParseInt(temp[1], 10, 64)
		case "packagenumberlength":
			myFileInfo.PackageNumberLength, err = strconv.ParseInt(temp[1], 10, 64)
		default:
			return nil
		}

		if err != nil {
			return nil
		}
	}
	return myFileInfo
}

type MyFileInfo struct {
	FileSize            int64
	FileName            string
	PackageCount        int64
	PackageNumberLength int64
}

//格式化对象为字串
func (this *MyFileInfo) ToString() string {
	var result string
	result = "title{:}" + "filesize=" + strconv.FormatInt(this.FileSize, 10) + "{:}" + "filename=" + this.FileName + "{:}" + "packagecount=" + strconv.FormatInt(this.PackageCount, 10) + "{:}" + "packagenumberlength=" + strconv.FormatInt(this.PackageNumberLength, 10)
	return result
}
