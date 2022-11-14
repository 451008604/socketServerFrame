package main

import (
	"fmt"
	"socketServerFrame/client/api"
	"socketServerFrame/client/base"
	"socketServerFrame/logs"
	"sync"
	"time"
)

func main() {
	logs.SetPrintMode(false)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for n := 0; n < 1; n++ {
		go func(n int) {
			conn := &base.CustomConnect{}
			conn.NewConnection("127.0.0.1", "7777")
			defer conn.SetBlocking()
			go func(n int) {
				i := 0
				for {
					i++
					msgId := uint32(2001)
					reqData := api.MarshalJsonData(api.PingReq{
						Msg:       fmt.Sprintf("C2S [connId %v][msgId %v]:ping -> %v", n, msgId, time.Now().UnixMilli()),
						TimeStamp: time.Now().UnixMilli(),
					})
					conn.SendMsg(msgId, []byte(reqData))
					logs.PrintLogInfoToConsole(reqData)
					time.Sleep(5 * time.Second)
				}
			}(n)
		}(n)
	}
	wg.Wait()
}
