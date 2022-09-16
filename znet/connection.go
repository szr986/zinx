package znet

import (
	"fmt"
	"net"

	"example.com/m/ziface"
)

type Connection struct {
	// 当前连接的socket tcp套接字
	Conn *net.TCPConn

	// 连接的id
	ConnID uint32

	// 当前连接的状态
	isClosed bool

	// 当前连接所绑定的API
	handleAPI ziface.HandleFunc

	// 告知当前连接已经退出停止的channel
	Exitchan chan bool
}

// 初始化连接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, callback_api ziface.HandleFunc) *Connection {
	c := &Connection{
		Conn:      conn,
		ConnID:    connID,
		handleAPI: callback_api,
		isClosed:  false,
		Exitchan:  make(chan bool, 1),
	}
	return c
}

// 连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine running...")
	defer fmt.Println("connID= ", c.ConnID, "Reader is exit,remote addr is ", c.Conn.RemoteAddr().String())
	defer c.Stop()

	for {
		// 读取客户端的数据到buf中,最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			continue
		}

		// 调用当前连接所绑定的handleAPI
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("ConnID", c.ConnID, "handle is error", err)
			break
		}

	}
}

// 启动连接，让当前连接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Start.. ConnID = ", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()
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
	// 回收资源
	close(c.Exitchan)

}

// 获取当前连接绑定的socket conn
func (c *Connection) GetIPConnection() *net.TCPConn {
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

// 发送数据,将数据发送给客户端
func (c *Connection) Send(data []byte) error {

}
