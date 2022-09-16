package main

import (
	"fmt"

	"example.com/m/ziface"
	"example.com/m/znet"
)

// 基于zinx开发的服务端应用程序

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// TEst PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// test handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte(" ping..."))
	if err != nil {
		fmt.Println("call back  ping error")
	}
}

// test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinxv0.2]")
	// 添加一个自定义router
	s.AddRouter(&PingRouter{})
	// 2. 启动server
	s.Serve()
}
