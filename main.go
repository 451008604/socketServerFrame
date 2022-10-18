package main

import (
	"socketServerFrame/client/api"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
	"socketServerFrame/znet"
	"time"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handler(req iface.IRequest) {
	// _, _ = req.GetConnection().GetTCPConnection().Write(req.GetData())

	reqData := api.MarshalJsonData(api.PingReq{
		Msg:       string(req.GetData()),
		TimeStamp: time.Now().UnixMilli(),
	})
	logs.PrintLogInfoToFile(reqData)
	req.GetConnection().SendMsg(1, []byte(reqData))
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Server()
}
