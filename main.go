package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/451008604/socketServerFrame/api"
	"github.com/451008604/socketServerFrame/iface"
	"github.com/451008604/socketServerFrame/logs"
	"github.com/451008604/socketServerFrame/network"
	pb "github.com/451008604/socketServerFrame/proto/bin"
)

func main() {
	// 捕获异常
	defer func() {
		if err := recover(); err != nil {
			logs.PrintLogPanic(fmt.Errorf("%v", err))
			// 阻塞防止主线程退出中断异步打印日志
			select {}
		}
	}()

	s := network.NewServer()
	go func(s iface.IServer) {
		for range time.Tick(time.Second * 3) {
			logs.PrintLogInfo(fmt.Sprint("当前线程数：", runtime.NumGoroutine(), "\t当前连接数量：", s.GetConnMgr().Len()))
		}
	}(s)
	s.SetOnConnStart(func(conn iface.IConnection) {
		conn.SetProperty("Client", conn.RemoteAddr())
	})
	s.AddRouter(uint32(pb.MessageID_PING), &api.PingRouter{})
	s.SetOnConnStop(func(conn iface.IConnection) {
		logs.PrintLogInfo(fmt.Sprintf("客户端%v下线", conn.GetProperty("Client")))
	})

	s.Server()
}
