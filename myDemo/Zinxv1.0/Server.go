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

// 创建连接之后执行的hook函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionBegin is called....")
	if err := conn.SendMsg(202, []byte("doConnectionBegin")); err != nil {
		fmt.Println("err")
	}

	// 设置连接属性
	fmt.Println("Set conn Name,Hoe```")
	conn.SetProperty("Name", "szr986")
	conn.SetProperty("github", "github.com/szr986")
}

// 连接断开前执行的hook函数
func DoConnectionStop(conn ziface.IConnection) {
	fmt.Println("=====> DoConnectionStop is called....")
	fmt.Println("connid = ", conn.GetConnID(), "is lost````")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("name = ", name)
	}

	if github, err := conn.GetProperty("github"); err == nil {
		fmt.Println("name = ", github)
	}
}

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinxV0.5]")
	// 注册连接的hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionStop)
	// 添加一个自定义router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 2. 启动server
	s.Serve()
}
