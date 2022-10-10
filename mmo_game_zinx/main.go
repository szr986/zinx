package main

import (
	"fmt"

	"example.com/m/mmo_game_zinx/apis"
	"example.com/m/mmo_game_zinx/core"
	"example.com/m/ziface"
	"example.com/m/znet"
)

func main() {
	// 创建server句柄
	s := znet.NewServer("MMO GAME ZINX")

	// 连接创建和销毁hook函数
	s.SetOnConnStart(OnConnectionAdd)
	// 注册路由
	s.AddRouter(2, &apis.WorldChatApi{})
	// 启动服务
	s.Serve()
}

func OnConnectionAdd(conn ziface.IConnection) {
	// 创建一个Player对象
	player := core.NewPlayer(conn)
	// 给客户端发送MsgID：1 的消息
	player.SyncPid()
	// 发送 200 广播消息 同步初始位置
	player.BroadCastStartPosition()
	core.WorldMgrObj.AddPlayer(player)
	conn.SetProperty("pid", player.Pid)
	fmt.Println("======> Player id = ", player.Pid)
}
