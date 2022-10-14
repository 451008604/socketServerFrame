package iface

/*
	IRequest 接口：
	实际上是把客户端请求的连接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface {
	GetConnection() // 获取请求连接信息
	GetData()       // 获取请求消息的数据
}
