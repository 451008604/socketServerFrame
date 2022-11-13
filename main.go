package main

import (
	"fmt"
	"runtime"
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

	// reqData := api.MarshalJsonData(api.PingReq{
	// 	Msg:       string(req.GetData()),
	// 	TimeStamp: time.Now().UnixMilli(),
	// })
	// if req.GetMsgID() == 2002 {
	// 	reqData = api.MarshalJsonData(api.PingReq{
	// 		Msg:       "MsgId 2002",
	// 		TimeStamp: time.Now().UnixMilli(),
	// 	})
	// }
	logs.PrintLogInfoToFile("服务成功接收心跳")
	// req.GetConnection().SendMsg(1, []byte(reqData))
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
	// s.AddRouter(2002, &PingRouter{})
	s.Server()
}
