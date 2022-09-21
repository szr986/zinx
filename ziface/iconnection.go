package ziface

import "net"

type IConnection interface {
	// 启动连接，让当前连接准备开始工作
	Start()
	// 停止连接，结束当前连接工作
	Stop()
	// 获取当前连接绑定的socket conn
	GetTCPConnection() *net.TCPConn
	// 获取当前连接模块的连接id
	GetConnID() uint32
	// 获取远程客户端的tcp状态 IP port
	RemoteAddr() net.Addr
	// 发送数据,将数据发送给客户端
	SendMsg(msgId uint32, data []byte) error
	// 设置连接属性
	SetProperty(key string, value interface{})
	// 获取连接属性
	GetProperty(key string) (interface{}, error)
	// 移除连接属性
	RemoveProperty(key string)
}

// 定义一个处理连接义务的方法(连接，数据内容，数据长度)
type HandleFunc func(*net.TCPConn, []byte, int) error
