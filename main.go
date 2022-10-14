package main

import (
	"socketServerFrame/iface"
	"socketServerFrame/znet"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handler(req iface.IRequest) {
	_, _ = req.GetConnection().GetTCPConnection().Write(req.GetData())
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Server()
}
