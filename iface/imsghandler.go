package iface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 异步处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
	StartWorkerPool()                       // 启动worker工作池
	SendMsgToTaskQueue(request IRequest)    // 将消息推入TaskQueue，等待Worker处理
}
