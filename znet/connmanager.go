package znet

import (
	"errors"
	"socketServerFrame/iface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]iface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                 // 连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]iface.IConnection),
	}
}

// Add 添加连接
func (c *ConnManager) Add(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	c.connections[conn.GetConnID()] = conn
}

// Remove 删除连接
func (c *ConnManager) Remove(conn iface.IConnection) {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	delete(c.connections, conn.GetConnID())
}

// Get 根据ConnID获取连接
func (c *ConnManager) Get(connID uint32) (iface.IConnection, error) {
	c.connLock.Lock()
	defer c.connLock.Unlock()
	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("未获取到绑定的conn")
	}
}

// Len 获取当前连接数量
func (c *ConnManager) Len() int {
	return len(c.connections)
}

// ClearConn 删除并停止所有连接
func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	defer c.connLock.Unlock()

	// 清理全部的connections信息
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
}
