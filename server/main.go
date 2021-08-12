package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"

	"coltcn.com/file-sync-server/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type File struct {
	Name    string
	Size    int64
	Content []byte
}

var conf config.Server

func main() {
	// 读取配置文件
	fmt.Println("读取配置文件")
	v := viper.New()
	v.SetConfigFile("config.yaml")
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置文件失败！"))
	}
	v.WatchConfig()

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件修改：", e.Name)
		if err := v.Unmarshal(&conf); err != nil {
			fmt.Println(err)
		}
	})

	if err := v.Unmarshal(&conf); err != nil {
		fmt.Println(err)
	}
	// 启动服务器
	fmt.Println("启动服务器...")
	listener, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		panic(fmt.Errorf("启动服务器失败：%s", err))
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("握手失败！", err)
		}
		go handlConn(conn)
	}
}

func handlConn(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr()
	fmt.Println("连接成功:", remoteAddr)

	//	buf := make([]byte, 1024)
	//	for {
	//		n, err := conn.Read(buf)
	//		if err != nil {
	//			fmt.Println("读取数据失败！", err)
	//			return
	//		}
	//		fmt.Printf("from %s data:%s\n", remoteAddr, string(buf[:n]))
	//	}

	for {
		// 先读取包长度
		lenByte := make([]byte, 4)
		if _, err := io.ReadFull(conn, lenByte); err != nil {
			if err == io.EOF {
				fmt.Println("客户端关闭连接")
				return
			}
			fmt.Println("读取包长度失败！", err)
			return
		}
		buff := bytes.NewBuffer(lenByte)
		var packageLen int32
		if err := binary.Read(buff, binary.LittleEndian, &packageLen); err != nil {
			fmt.Println("读取包长度失败！", err)
			return
		}
		//packageLen := binary.BigEndian.Uint32(lenByte)
		fmt.Printf("包大小：%d\n", packageLen)
		content := make([]byte, packageLen)
		io.ReadFull(conn, content)
		fmt.Printf("%d-%d", lenByte[:], content[:4])
		file, err := readData(content)
		if err != nil {
			return
		}
		fmt.Printf("接收文件：%s\n", file.Name)

	}
}

func readData(data []byte) (*File, error) {
	file := new(File)
	dc := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dc.Decode(&file); err != nil {
		fmt.Println("解析数据失败！")
		return nil, err
	}
	return file, nil
}
