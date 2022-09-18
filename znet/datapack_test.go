package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 封包拆包测试
func TestDataPack(t *testing.T) {
	// 模拟服务器
	// 1. 创建socketTCP
	listener, err := net.Listen("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("server listener err:", err)
		return
	}
	// 创建一个go从客户端处理
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			go func(conn net.Conn) {
				// 处理客户端请求
				// 定义一个拆包对象
				dp := NewDataPack()
				for {
					// 1. 第一次读head
					headData := make([]byte, dp.GetHeadLen())
					// readfull 读满切片内存
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err:", err)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack err:", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						// msg有数据，需要进行第二次读取
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						// 根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("read data err:", err)
							return
						}
						// 完整消息已经读完
						fmt.Println("---> Received MsgID", msg.Id,
							",datalen = ", msg.DataLen,
							"data = ", string(msg.Data))
					}
				}

			}(conn)
		}
	}()
	// 2. 从客户端读取数据进行拆包处理

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7778")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}

	// 创建一个封包对象
	dp := NewDataPack()
	// 模拟粘包,封装两个Message一同发送
	// 封装第一个Msg1
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
	}
	// 封装Msg2
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', 'a', 'a'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
	}
	// 将两个包黏在一起
	sendData1 = append(sendData1, sendData2...)
	// 一次性发送给服务端
	conn.Write(sendData1)

	// 客户端阻塞
	select {}
}
