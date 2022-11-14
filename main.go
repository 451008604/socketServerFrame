package main

import (
	"fmt"
	"runtime"
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
	t := time.Now().UnixMilli()
	pingReq := api.PingReq{}
	api.UnmarshalJsonData(req.GetData(), &pingReq)

	resData := api.MarshalJsonData(api.PingRes{
		Msg: t - pingReq.TimeStamp,
	})
	req.GetConnection().SendMsg(1, []byte(resData))
}

func main() {
	s := znet.NewServer()
	go func(s iface.IServer) {
		for range time.Tick(time.Second * 3) {
			logs.PrintLogInfoToConsole(fmt.Sprint("当前线程数：", runtime.NumGoroutine(), "\t当前连接数量：", s.GetConnMgr().Len()))
		}
	}(s)
	s.SetOnConnStart(func(conn iface.IConnection) {
		conn.SendBuffMsg(1, []byte("连接开始"))
	})
	s.SetOnConnStop(func(conn iface.IConnection) {
		conn.SendBuffMsg(2, []byte("连接关闭"))
	})
	s.AddRouter(2001, &PingRouter{})
	s.Server()
}
