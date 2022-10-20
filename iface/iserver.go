package iface

// IServer 定义服务器接口
type IServer interface {
	// Start 启动服务器
	Start()
	// Stop 停止服务器
	Stop()
	// Server 开启业务服务
	Server()
	// AddRouter 路由功能：给当前服务注册一个路由业务方法，提供给客户端连接使用
	AddRouter(msgId uint32, router IRouter)
	// GetConnMgr 获取连接管理器
	GetConnMgr() IConnManager
}
