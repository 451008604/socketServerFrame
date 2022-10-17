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
	for n := 0; n < 5; n++ {
		go func(n int) {
			conn := &base.CustomConnect{}
			conn.NewConnection("127.0.0.1", "7777")
			defer conn.SetBlocking()
			go func(n int) {
				i := 0
				for true {
					i++
					reqData := api.MarshalJsonData(api.PingReq{
						Msg:       fmt.Sprintf("[connId %v]:ping -> %v", n, time.Now().UnixMilli()),
						TimeStamp: time.Now().UnixMilli(),
					})
					conn.SendMsg(2001, []byte(reqData))
					time.Sleep(1 * time.Second)
				}
			}(n)
		}(n)
	}
	wg.Wait()
}
