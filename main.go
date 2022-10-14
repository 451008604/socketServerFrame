package main

import (
	"SocketServerFrame/znet"
	"socketServerFrame/iface"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handler(req iface.IRequest) {
	_, _ = req.GetConnection().GetTCPConnection().Write(req.GetData())
}

func main() {
	s := znet.NewServer("myServer")
	s.AddRouter(&PingRouter{})
	s.Server()
}
