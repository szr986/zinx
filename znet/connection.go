package znet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"example.com/m/utils"
	"example.com/m/ziface"
)

type Connection struct {
	// 当前连接的socket tcp套接字
	Conn *net.TCPConn

	// 连接的id
	ConnID uint32

	// 当前连接的状态
	isClosed bool

	// 告知当前连接已经退出停止的channel
	Exitchan chan bool

	// 无缓冲的管道，用于读写goroutine之间的消息通信
	msgChan chan []byte

	// 消息的管理MsgID，和对应的处理业务API
	MsgHandle ziface.IMsgHandler
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandle,
		msgChan:   make(chan []byte),
		isClosed:  false,
		Exitchan:  make(chan bool, 1),
	}
	return c
}

// 写消息的goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWriter() {
	fmt.Println("[Wriiter goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), "[conn Wriiter exit!]")

	// 不断的阻塞等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err", err)
				return
			}
		case <-c.Exitchan:
			// 代表Reader已经退出，告知writer退出
			return
		}
	}
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Goroutine running...]")
	defer fmt.Println("connID= ", c.ConnID, "Reader is exit,remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中,最大512字节
		// buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		// _, err := c.Conn.Read(buf)
		// if err != nil {
		// 	fmt.Println("recv buf err", err)
		// 	continue
		// }

		// 创建一个拆包解包的对象
		dp := NewDataPack()
		// 读取客户端的MsgHead 二进制流 8字节，
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("recv msg head err", err)
			break
		}

		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack msg err", err)
			break
		}
		// 拆包，得到MsgID 和 MsgDataLen，放在一个Msg对象中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err", err)
				break
			}
		}
		msg.SetData(data)
		// 根据datalen 再次读取Data，放在msg.data

		// 得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		// 判断是否开启工作池
		if utils.GlobalObject.WorkerPoolSize > 0 {
			// 将消息发给worker工作池
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从路由中，找到注册绑定的conn对应的router调用
			// 根据绑定好的msgid 找到对应的处理API
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}
}

// 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start.. ConnID = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
	go c.StartWriter()
	//TODO 启动从当前连接写数据的业务

}

// 停止连接，结束当前连接工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop..ConnID = ", c.ConnID)

	if c.isClosed == true {
		return
	}
	c.isClosed = true
	// 关闭socket连接
	c.Conn.Close()
	// 告知writer关闭
	c.Exitchan <- true

	// 回收资源
	close(c.Exitchan)
	close(c.msgChan)

}

// 获取当前连接绑定的socket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的连接id
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的tcp状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法,将我们要发送给客户端的数据进行封包再发送
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("connection is closed when send msg!")
	}

	// 将data进行封包,MsgdataLen,Msgid,MsgData
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack msg error id= ", msgId)
		return errors.New("pack msg err")
	}
	// 将数据写回管道
	c.msgChan <- binaryMsg

	return nil
}
