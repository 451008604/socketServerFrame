package iface

// IMsgHandler 消息管理抽象层
type IMsgHandler interface {
	DoMsgHandler(request IRequest)          // 异步处理消息
	AddRouter(msgId uint32, router IRouter) // 为消息添加具体的处理逻辑
}
