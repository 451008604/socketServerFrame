package main

import (
	"fmt"
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
					conn.SendMsg([]byte(fmt.Sprintf("[connId %v]:send -> %v\n", n, i)))
					time.Sleep(1 * time.Second)
				}
			}(n)
		}(n)
	}
	wg.Wait()
}
