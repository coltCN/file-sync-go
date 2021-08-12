package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net"
	"path/filepath"
)

type File struct {
	Name    string
	Size    int64
	Content []byte
}

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(fmt.Errorf("无法连接服务器！"))
	}
	defer client.Close()
	root := "/mnt/c/Users/BZL/Pictures/背景"
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 跳过文件夹
		if info.IsDir() {
			return nil
		}
		// fmt.Printf("路径：%s -> 名称：%s Size: %d \n", path, info.Name(), info.Size())
		sendFile(client, path, info)
		return nil
	})
}

/// 使用 gob 处理二进制，只适用于客户端服务端都使用gob包进行编码和解码的情况
func sendFile(client net.Conn, path string, info fs.FileInfo) {
	//data := fmt.Sprintf("名称：%s 大小：%d", (*info).Name(), (*info).Size())
	//client.Write([]byte(data))

	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败！", err)
		return
	}
	file := File{info.Name(), info.Size(), content}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(file); err != nil {
		fmt.Println("消息转换失败！", err)
		return
	}
	dataPackage := new(bytes.Buffer)
	// 写入长度 并转换为4个字节长度，即int32
	packageLen := int32(buf.Len())
	binary.Write(dataPackage, binary.LittleEndian, packageLen)
	binary.Write(dataPackage, binary.LittleEndian, buf.Bytes())
	fmt.Printf("文件：%s,包长：%d,发送数据大小：%d\n", info.Name(), packageLen, dataPackage.Len())
	client.Write(dataPackage.Bytes())
}
