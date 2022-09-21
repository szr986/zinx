package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"example.com/m/znet"
)

// 模拟客户端
func main() {
	fmt.Println("client1 start..")
	time.Sleep(1 * time.Second)
	// 1.直接连接远程服务器，得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err,exit!", err)
		return
	}
	// 2. 连接调用write，写数据
	for {
		// 发送封包的message消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(1, []byte("ZinxV0.6 client1 Test Msg")))
		if err != nil {
			fmt.Println("Pack err,", err)
			return
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write err,", err)
			return
		}
		// 服务器回复一个message数据
		// 先读取流中的head部分
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head err", err)
			break
		}
		// 将二进制的head拆包到msg结构体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack msghead err,", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			// 再根据datalen进行第二次读取,将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msgHead.GetMsgLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msgdata err,", err)
				return
			}

			fmt.Println("----> Received Server Msg :ID = ", msg.Id,
				"len = ", msg.DataLen, ",data= ", string(msg.Data))
		}

		// cpu阻塞
		time.Sleep(1 * time.Second)
	}

}
