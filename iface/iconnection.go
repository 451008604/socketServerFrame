package iface

import "net"

type IConnection interface {
	Start()                         // 启动连接
	Stop()                          // 停止连接
	GetTCPConnection() *net.TCPConn // 从当前连接获取原始的Socket TCPConn
	GetConnID() uint32              // 获取当前连接ID
	RemoteAddr() net.Addr           // 获取客户端地址信息
}

// HandFunc 统一处理连接业务的接口
type HandFunc func(conn *net.TCPConn, reqMsgData []byte, dataLength int) error
