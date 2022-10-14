package znet

import "socketServerFrame/iface"

// BaseRouter router基类，业务逻辑根据需要对基类方法重写
type BaseRouter struct{}

func (b *BaseRouter) PreHandler(req iface.IRequest) {
}

func (b *BaseRouter) Handler(req iface.IRequest) {
}

func (b *BaseRouter) AfterHandler(req iface.IRequest) {
}
