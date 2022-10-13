package core

import (
	"fmt"
	"sync"

	"example.com/m/mmo_game_zinx/db"
)

// AOI地图中的一个格子类型

type Grid struct {
	// 格子ID
	GID int
	// 格子左边界坐标
	MinX int
	// 格子右边界坐标
	MaxX int
	// 格子上边界坐标
	MinY int
	// 格子下边界坐标
	MaxY int
	// 当前格子内玩家或物体的ID集合
	playerIDs map[int]bool
	// 保护当前集合的锁
	pIDLock sync.RWMutex
}

// 初始化当前格子的方法
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// 给格子添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true

	rdb := db.GetRedisClient()
	Istrue, err := rdb.SIsMember("grid:"+string(g.GID), playerID).Result()
	if err != nil {
		fmt.Println("get from redis err:", err)
		return
	}
	if Istrue == true {
		fmt.Println("player already exists in grid : ", g.GID)
	}

	rdb.SAdd("grid:"+string(g.GID), playerID)
	fmt.Println("Redis Grid : ", g.GID, "Player : ", rdb.SMembers("grid:"+string(g.GID)).String())
}

// 从格子删除一个玩家
func (g *Grid) Remove(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// 得到当前格子的所有玩家
func (g *Grid) GetPlayerIDs() (playerIDs []int) {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()

	for k, _ := range g.playerIDs {
		playerIDs = append(playerIDs, k)
	}

	return
}

// 调试使用，打印格子信息
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id:%d,minX:%d,maxX:%d,minY:%d,maxY:%d,playerIDs:%v", g.GID,
		g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
