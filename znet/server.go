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
	Router ziface.IRouter
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

			// 将处理新连接的业务方法和connection进行绑定，得到我们的连接模块
			dealConn := NewConnection(conn, cid, s.Router)
			cid++
			// 启动业务
			go dealConn.Start()
		}
	}()

}

func (s *Server) Stop() {
	// TODO 将服务器的资源，状态或者一些已经开辟的连接信息进行停止或回收
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()
	//TODO 做一些额外业务

	// 阻塞状态
	select {}
}

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router succ")
}

// 初始化server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPversion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      8999,
		Router:    nil,
	}

	return s
}
