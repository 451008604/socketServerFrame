package iface

import "net"

type IConnection interface {
	Start() // 启动连接
	Stop()  // 停止连接

	GetTCPConnection() *net.TCPConn // 从当前连接获取原始的Socket TCPConn
	GetConnID() uint32              // 获取当前连接ID
	RemoteAddr() net.Addr           // 获取客户端地址信息

	SendMsg(msgId uint32, data []byte)     // 发送消息给客户端（无缓冲）
	SendBuffMsg(msgId uint32, data []byte) // 发送消息给客户端（有缓冲）

	SetProperty(key string, value interface{})  // 设置连接属性
	GetProperty(key string) (value interface{}) // 获取连接属性
	RemoveProperty(key string)                  // 删除连接属性
}

// HandFunc 统一处理连接业务的接口
type HandFunc func(conn *net.TCPConn, reqMsgData []byte, dataLength int) error
