package main

import (
	"fmt"
	"runtime"
	"socketServerFrame/api"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
	pb "socketServerFrame/proto/bin"
	"socketServerFrame/znet"
	"time"
)

func main() {
	s := znet.NewServer()
	go func(s iface.IServer) {
		for range time.Tick(time.Second * 3) {
			logs.PrintLogInfoToConsole(fmt.Sprint("当前线程数：", runtime.NumGoroutine(), "\t当前连接数量：", s.GetConnMgr().Len()))
		}
	}(s)
	s.SetOnConnStart(func(conn iface.IConnection) {
		conn.SetProperty("Client", conn.RemoteAddr())
	})
	s.AddRouter(uint32(pb.MessageID_PING), &api.PingRouter{})
	s.SetOnConnStop(func(conn iface.IConnection) {
		logs.PrintLogInfoToConsole(fmt.Sprintf("客户端%v下线", conn.GetProperty("Client")))
	})
	s.Server()
}
