package main

import (
	"google.golang.org/protobuf/proto"
	"socketServerFrame/client/base"
	"socketServerFrame/logs"
	pb "socketServerFrame/proto/bin"
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

					data := &pb.Ping{TimeStamp: time.Now().UnixMicro()}
					marshal, err := proto.Marshal(data)
					if err != nil {
						return
					}
					conn.SendMsg(uint32(pb.MessageID_PING), marshal)
					logs.PrintLogInfoToConsole(data.String())
					time.Sleep(5 * time.Second)
				}
			}(n)
		}(n)
	}
	wg.Wait()
}
