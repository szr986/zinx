package apis

import (
	"fmt"

	"example.com/m/mmo_game_zinx/core"
	"example.com/m/mmo_game_zinx/pb"
	"example.com/m/ziface"
	"example.com/m/znet"
	"google.golang.org/protobuf/proto"
)

// 玩家移动

type MoveApi struct {
	znet.BaseRouter
}

func (m *MoveApi) Handle(request ziface.IRequest) {
	proto_msg := &pb.Position{}
	err := proto.Unmarshal(request.GetData(), proto_msg)
	if err != nil {
		fmt.Println("Move Position unmarshal error:", err)
		return
	}

	pid, err := request.GetConnection().GetProperty("pid")
	if err != nil {
		fmt.Println("Get mmove pid error:", err)
		return
	}

	fmt.Printf("Player pid:%d,move:(%f,%f,%f,%f)\n", pid, proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)

	player := core.WorldMgrObj.GetPlayerByPid(pid.(int32))
	// 广播更新当前玩家坐标
	player.UpdatePos(proto_msg.X, proto_msg.Y, proto_msg.Z, proto_msg.V)
}
