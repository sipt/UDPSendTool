package io

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	DATA_LENGTH = 1024
)

type FileOp struct {
	ReadData  chan []byte
	WriteData chan []byte
}

/**
*文件读取
 */
func (this *FileOp) ReadFile(fileName string) {
	//打开文件
	fi, err := os.Open(fileName)
	if err != nil {
		this.ReadData <- nil
		panic(err)
	}
	//函数结束关闭文件
	defer fi.Close()

	//定义读取流
	r := bufio.NewReader(fi)
	//当前slice里数据长度
	var n int
	for {
		//引用传递，所以每次创建一个新的
		buf := make([]byte, DATA_LENGTH)
		n, err = r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 { //数据长度为0时，退出函数
			this.ReadData <- nil
			return
		} else { //将数据放入channel
			this.ReadData <- buf[:n]
		}
	}
}

/**
*文件写入
 */
func (this *FileOp) WriteFile(fileName string) {
	_, err := os.Stat(fileName)
	//查看文件名是否存在
	if err == nil || os.IsExist(err) {
		//err = os.Remove(fileName)
		fmt.Println("文件已存在！")
		return
	}
	//创建文件
	fi, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	//函数结束关闭文件
	defer fi.Close()

	//定义写入流
	w := bufio.NewWriter(fi)

	for {
		//引用传递，所以每次创建一个新的
		buf := <-this.WriteData
		if buf != nil && len(buf) != 0 {
			w.Write(buf)
			w.Flush()
		} else {
			fmt.Println("写入结束")
			return
		}
	}
}
