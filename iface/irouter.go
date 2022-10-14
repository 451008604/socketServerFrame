package iface

/*
	路由接口， 这里面路由是 使用框架者给该连接自定的 处理业务方法
	路由里的IRequest 则包含用该连接的连接信息和该连接的请求数据信息
*/
type IRouter interface {
	PreHandler(req IRequest)   // 处理conn业务之前
	Handler(req IRequest)      // 处理conn业务逻辑
	AfterHandler(req IRequest) // 处理conn业务之后
}
