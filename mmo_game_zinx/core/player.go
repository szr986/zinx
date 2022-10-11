package core

import (
	"fmt"
	"math/rand"
	"sync"

	"example.com/m/mmo_game_zinx/pb"
	"example.com/m/ziface"
	"google.golang.org/protobuf/proto"
)

// 玩家对象

type Player struct {
	Pid  int32              //玩家ID
	Conn ziface.IConnection //当前玩家连接
	X    float32
	Y    float32
	Z    float32
	V    float32
}

// PID生成器
var PidGen int32 = 1 // ID计数器
var PidLock sync.Mutex

// 创建一个玩家的方法
func NewPlayer(conn ziface.IConnection) *Player {
	// 生成一个玩家ID
	PidLock.Lock()
	id := PidGen
	PidGen++
	PidLock.Unlock()
	// 创建一个玩家对象
	p := &Player{
		Pid:  id,
		Conn: conn,
		X:    float32(100 + rand.Intn(10)),
		Y:    0,
		Z:    float32(100 + rand.Intn(10)),
		V:    0, // 角度为0
	}
	return p
}

// 提供一个发送给客户端消息的方法
func (p *Player) SendMsg(msgId uint32, data proto.Message) {
	// 将proto Message结构体序列化，转化成二进制
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err:", err)
		return
	}
	// 将二进制信息发送
	if p.Conn == nil {
		fmt.Println("conn in Player is nil")
		return
	}

	if err := p.Conn.SendMsg(msgId, msg); err != nil {
		fmt.Println("Player Send msg err:", err)
		return
	}

	return
}

// 告知客户端挖按键Pid
func (p *Player) SyncPid() {
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	p.SendMsg(1, proto_msg)
}

// 广播玩家的出生地点
func (p *Player) BroadCastStartPosition() {
	// 组件MsgID 200 的proto协议
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	p.SendMsg(200, proto_msg)
}

// 玩家广播世界聊天消息
func (p *Player) Talk(content string) {
	// 组件MsgId=200的 proto数据
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  1,
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}
	// 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()
	// 发送消息
	for _, player := range players {
		// 分别给对应的客户端发送消息
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) SyncSurrounding() {
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)
	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}

	// 将周围玩家数据发给自己
	players_proto_msg := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			Pid: player.Pid,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		players_proto_msg = append(players_proto_msg, p)
	}

	SyncPlayer_proto_msg := &pb.SyncPlayer{
		Ps: players_proto_msg[:],
	}

	p.SendMsg(202, SyncPlayer_proto_msg)
}

// 广播当前玩家位置移动信息

func (p *Player) UpdatePos(x, y, z, v float32) {
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v
	proto_msg := &pb.BroadCast{
		Pid: p.Pid,
		Tp:  4,
		Data: &pb.BroadCast_P{
			P: &pb.Position{X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	players := p.GetSurroudingPlayers()

	for _, player := range players {
		player.SendMsg(200, proto_msg)
	}
}

func (p *Player) GetSurroudingPlayers() []*Player {
	pids := WorldMgrObj.AoiMgr.GetPidsByPos(p.X, p.Z)

	players := make([]*Player, 0, len(pids))
	for _, pid := range pids {
		players = append(players, WorldMgrObj.GetPlayerByPid(int32(pid)))
	}

	return players
}

func (p *Player) Offline() {
	players := p.GetSurroudingPlayers()
	proto_msg := &pb.SyncPid{
		Pid: p.Pid,
	}

	for _, player := range players {
		player.SendMsg(201, proto_msg)
	}

	WorldMgrObj.AoiMgr.RemoveFromGridByPos(int(p.Pid), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPid(p.Pid)
}
