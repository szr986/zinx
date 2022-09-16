package main

import "example.com/m/znet"

// 基于zinx开发的服务端应用程序

func main() {
	// 1.创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinxv0.2]")
	// 2. 启动server
	s.Serve()
}
