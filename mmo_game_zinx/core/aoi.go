package core

import "fmt"

// 定义一些AOI边界值
const (
	AOI_MIN_X  int = 0
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 0
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

// AOI 区域管理模块

type AOIManager struct {
	// 区域的左边界坐标
	MinX int
	// 区域的右边界坐标
	MaxX int
	// X方向格子的数量
	CntsX int
	// 区域的上边界坐标
	MinY int
	// 区域的下边界坐标
	MaxY int
	// Y方向的格子数量
	CntsY int
	// 当前区域中有哪些格子
	grids map[int]*Grid
}

// 初始化AOI区域管理器
func NewAOIManager(minX, maxX, cntsX, minY, maxY, cntsY int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:  minX,
		MaxX:  maxX,
		CntsX: cntsX,
		MinY:  minY,
		MaxY:  maxY,
		CntsY: cntsY,
		grids: make(map[int]*Grid),
	}

	// 给AOI初始化区域的格子进行编号
	for y := 0; y < cntsY; y++ {
		for x := 0; x < cntsX; x++ {
			// 根据X,Y编号
			gid := y*cntsX + x
			// 初始化gid格子
			aoiMgr.grids[gid] = NewGrid(gid, aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}

	return aoiMgr
}

// 得到每个格子在X轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return ((m.MaxX - m.MinX) / m.CntsX)
}

// 得到每个格子在Y轴方向的高度
func (m *AOIManager) gridLength() int {
	return ((m.MaxY - m.MinY) / m.CntsY)
}

// 打印格子信息
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManager:\n MinX:%d,MinY:%d,MaxX:%d,MaxY:%d,cntsX:%d,cntsY:%d \n,grids:",
		m.MinX, m.MinY, m.MaxX, m.MaxY, m.CntsX, m.CntsY)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// 根据格子gid得到周边九宫格格子的集合
func (m *AOIManager) GetSurroudGridsByGid(gID int) (grids []*Grid) {
	// 判断gid是否在AOIManager中
	if _, ok := m.grids[gID]; !ok {
		return
	}
	// 初始化grids返回值切片
	grids = append(grids, m.grids[gID])
	// gid的左边是否有格子，右边是否有格子
	// 需要通过gid格子得到当前格子x轴的编号- idx = id % nx
	// 判断idx坐标左边是否还有格子，有的话放入gidsX中
	// 是否右边还有格子 ，放入gidsX中
	idx := gID % m.CntsX
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	if idx < m.CntsX-1 {
		grids = append(grids, m.grids[gID+1])
	}
	// 遍历gidsX 集合，判断上下
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	for _, v := range gidsX {
		idy := v / m.CntsY
		if idy > 0 {
			grids = append(grids, m.grids[v-m.CntsX])
		}
		if idy < m.CntsY-1 {
			grids = append(grids, m.grids[v+m.CntsX])
		}
	}
	return
}

func (m *AOIManager) GetGidByPos(x, y float32) int {
	idx := (int(x) - m.MinX) / m.gridWidth()
	idy := (int(y) - m.MinY) / m.gridLength()

	return idy*m.CntsX + idx
}

// 通过横纵坐标得到周边九宫格全部playerIDs
func (m *AOIManager) GetPidsByPos(x, y float32) (playerIDs []int) {
	// 得到当前玩家的GID格子id
	gID := m.GetGidByPos(x, y)

	grids := m.GetSurroudGridsByGid(gID)

	for _, v := range grids {
		playerIDs = append(playerIDs, v.GetPlayerIDs()...)
		fmt.Println(" ====>", v.playerIDs, v.GetPlayerIDs())
	}
	return
}

// 添加一个player到格子中
func (m *AOIManager) AddPidToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// 从格子中移除一个player
func (m *AOIManager) RemovePidFromGrid(pID, gID int) {
	m.grids[gID].Remove(pID)
}

// 通过GID获取全部的playerID
func (m *AOIManager) GetPidsByGid(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// 通过坐标将Player添加到格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// 通过坐标从格子中删除一个player
func (m *AOIManager) RemoveFromGridByPos(pID int, x, y float32) {
	gID := m.GetGidByPos(x, y)
	grid := m.grids[gID]
	grid.Remove(pID)
}
