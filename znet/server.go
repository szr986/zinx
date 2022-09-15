package znet

import (
	"fmt"
	"net"

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
}

func (s *Server) Start() {
	fmt.Printf("[Start]Server Listenner at IP:%s,Port %d,is staring/n", s.IP, s.Port)
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
		// 3. 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端连接过来，阻塞会返回
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 客户端已经建立连接，conn是对应的客户端连接句柄，做一个最基本的最大512字节长度的回显业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err", err)
						continue
					}

					fmt.Printf("recv client buf %s,cnt %d\n", buf, cnt)
					// 回显
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write back buf err", err)
						continue
					}
				}
			}()
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

// 初始化server模块的方法
func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPversion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
