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
	fmt.Println("Call Router Handle...")
	// 先读取客户端的数据
	fmt.Println("recv from client msgId=", request.GetMsgId(),
		"data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println("err")
	}
}

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinxV0.5]")
	// 添加一个自定义router
	s.AddRouter(&PingRouter{})
	// 2. 启动server
	s.Serve()
}
