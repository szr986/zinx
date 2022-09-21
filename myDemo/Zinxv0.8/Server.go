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

// test handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	// 先读取客户端的数据
	fmt.Println("recv from client msgId=", request.GetMsgId(),
		"data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("err")
	}
}

// hellozinx test 自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// test handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloRouter Handle...")
	// 先读取客户端的数据
	fmt.Println("recv from client msgId=", request.GetMsgId(),
		"data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello"))
	if err != nil {
		fmt.Println("err")
	}
}

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinxV0.5]")
	// 添加一个自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})
	// 2. 启动server
	s.Serve()
}
