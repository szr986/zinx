package apis

import (
	"fmt"

	"example.com/m/mmo_game_zinx/core"
	"example.com/m/mmo_game_zinx/pb"
	"example.com/m/ziface"
	"example.com/m/znet"
	"google.golang.org/protobuf/proto"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(request ziface.IRequest) {
	// 解析客户端传递进来的proto协议
	proto_msg := &pb.Talk{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("talk unmarshal error:", err)
		return
	}
	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("world chat get pid err:", err)
		return
	}
	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	player.Talk(proto_msg.Content)
}
