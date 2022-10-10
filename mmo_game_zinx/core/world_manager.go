package core

import "sync"

type WorldManager struct {
	AoiMgr *AOIManager

	Players map[int32]*Player

	pLock sync.RWMutex
}

// 提供一个对外的世界管理模块句柄（全局
var WorldMgrObj *WorldManager

// 初始化方法
func init() {
	WorldMgrObj = &WorldManager{
		// 创建世界AOI地图
		AoiMgr: NewAOIManager(AOI_MIN_X, AOI_MAX_X, AOI_CNTS_X, AOI_MIN_Y, AOI_MAX_Y, AOI_CNTS_Y),
		// 初始化Players集合
		Players: make(map[int32]*Player),
	}
}

// 添加一个玩家
func (wm *WorldManager) AddPlayer(player *Player) {
	wm.pLock.Lock()
	wm.Players[player.Pid] = player
	wm.pLock.Unlock()
	// 将player添加到AOIMANAGER中
	wm.AoiMgr.AddToGridByPos(int(player.Pid), player.X, player.Y)
}

// 删除一个玩家
func (wm *WorldManager) RemovePlayerByPid(pid int32) {
	player := wm.Players[pid]
	wm.AoiMgr.RemoveFromGridByPos(int(pid), player.X, player.Y)

	wm.pLock.Lock()
	delete(wm.Players, pid)
	wm.pLock.Unlock()
}

// 通过玩家ID查询Player对象
func (wm *WorldManager) GetPlayerByPid(pid int32) *Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()
	return wm.Players[pid]
}

// 获取全部的在线玩家
func (wm *WorldManager) GetAllPlayers() []*Player {
	wm.pLock.RLock()
	defer wm.pLock.RUnlock()

	players := make([]*Player, 0)

	for _, v := range wm.Players {
		players = append(players, v)
	}

	return players
}
