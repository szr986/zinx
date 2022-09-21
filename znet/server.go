package znet

import (
	"fmt"
	"net"

	"example.com/m/utils"
	"example.com/m/ziface"
)

// iServer的接口实现，定义一个Server的服务器模块
type Server struct {
	// 服务器名称
	Name string
	// 服务器绑定的ip版本
	IPversion string
	// 服务器监听IP
	IP string
	// 服务器监听端口
	Port int
	// 当前Server添加一个router，server注册的连接对应的处理业务
	MsgHandler ziface.IMsgHandler
	// server的连接管理器
	ConnMgr ziface.IConnManager
	// 创建连接之后的hook函数
	OnConnStart func(conn ziface.IConnection)
	// 销毁连接之前后的hook函数
	OnConnStop func(conn ziface.IConnection)
}

func (s *Server) Start() {
	fmt.Printf("[Zinx]Server Name:%s,Listenner at IP:%s,Port:%d is staring\n",
		utils.GlobalObject.Name,
		utils.GlobalObject.Host,
		utils.GlobalObject.TcpPort)
	fmt.Printf("[Zinx]Version %s,Maxconn %d,MaxPackageSize: %d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.Maxconn,
		utils.GlobalObject.MaxPackageSize)
	// 1.获取一个TCP的Addr
	go func() {
		// 开启消息队列及worker工作池
		s.MsgHandler.StartWorkerPool()

		addr, err := net.ResolveTCPAddr(s.IPversion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		// 2.监听服务器的地址
		listenner, err := net.ListenTCP(s.IPversion, addr)
		if err != nil {
			fmt.Println("listen", s.IPversion, "err", err)
			return
		}

		fmt.Println("start Zinx server success", s.Name, "succ,Listenning..")
		var cid uint32
		cid = 0

		// 3. 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			// 设置最大连接个数的判断，如果超过最大连接的数量，则关闭此连接
			if s.ConnMgr.Len() >= utils.GlobalObject.Maxconn {
				// 给客户端响应一个超出最大连接的错误包
				fmt.Println("TOO MANY connections! Maxconn = ", utils.GlobalObject.Maxconn)
				conn.Close()
				continue
			}

			// 将处理新连接的业务方法和connection进行绑定，得到我们的连接模块
			dealConn := NewConnection(s, conn, cid, s.MsgHandler)
			cid++
			// 启动业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	//将服务器的资源，状态或者一些已经开辟的连接信息进行停止或回收
	fmt.Println("[STOP] zinx server name ", s.Name)
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	//TODO 做一些额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add router succ")
}

// 初始化server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPversion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       8999,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}

	return s
}

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart 钩子函数
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call on OnConnStart()...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call on OnConnStop()...")
		s.OnConnStop(conn)
	}
}
