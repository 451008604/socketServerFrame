package main

import (
	"fmt"
	"socketServerFrame/client/api"
	"socketServerFrame/client/base"
	"sync"
	"time"
)

func main() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for n := 0; n < 1; n++ {
		go func(n int) {
			conn := &base.CustomConnect{}
			conn.NewConnection("127.0.0.1", "7777")
			defer conn.SetBlocking()
			go func(n int) {
				i := 0
				for true {
					i++
					msgId := uint32(2001)
					reqData := api.MarshalJsonData(api.PingReq{
						Msg:       fmt.Sprintf("[connId %v][msgId %v]:ping -> %v", n, msgId, time.Now().UnixMilli()),
						TimeStamp: time.Now().UnixMilli(),
					})
					conn.SendMsg(msgId, []byte(reqData))
					msgId = 2002
					conn.SendMsg(msgId, []byte(reqData))
					time.Sleep(1 * time.Second)
				}
			}(n)
		}(n)
	}
	wg.Wait()
}
