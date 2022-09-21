package znet

import (
	"errors"
	"fmt"
	"sync"

	"example.com/m/ziface"
)

// 链接管理模块

type ConnManager struct {
	connections map[uint32]ziface.IConnection //管理的连接信息集合
	connLock    sync.RWMutex                  //读写锁
}

// 创建连接
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

// 添加连接
func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 将conn加入到connManager中
	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connID=,", conn.GetConnID(), "connection add to ConnManager succ,conn num = ", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connID=,", conn.GetConnID(), "connection remove from ConnManager succ,conn num = ", connMgr.Len())

}

// 根据connID获取连接
func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源，加写锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		// 找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not found!")
	}
}

// 得到当前连接总数
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除所有链接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	// 删除conn并停止conn工作
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear all connections success!")
}
