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
			conn := base.NewConnection("127.0.0.1", "7777")
			defer conn.SetBlocking()
			fmt.Println(fmt.Sprintf("connect Id %v", n))

			go func() {
				i := 0
				for true {
					i++
					conn.SendMsg([]byte(fmt.Sprintf("%v\t", i)))
					time.Sleep(1 * time.Second)
				}
			}()
		}(n)
	}
	wg.Wait()
}
