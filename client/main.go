package main

import (
	"SocketServerFrame/client/base"
	"fmt"
	"time"
)

func main() {
	base.NewConnection("127.0.0.1", "7777")

	// for true {
	base.SendMsg([]byte("ping test"))
	// 	time.Sleep(1 * time.Second)
	// }

	go func() {
		i := 0
		for true {
			i++
			base.SendMsg([]byte(fmt.Sprintf("%v\n", i)))
			time.Sleep(1 * time.Second)
		}
	}()

	base.SetBlocking()
}
