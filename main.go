package main

import "SocketServerFrame/znet"

func main() {
	s := znet.NewServer("myServer")
	s.Server()
}
