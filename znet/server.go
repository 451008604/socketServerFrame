package znet

import (
	"SocketServerFrame/ziface"
	"fmt"
	"net"
	"time"
)

// Server 定义Server服务类实现IServer接口
type Server struct {
	Name      string // 服务器名称
	IPVersion string // tcp4 or other
	IP        string // IP地址
	Port      int    // 服务端口
}

func (s *Server) Start() {
	// 开启一个go去做服务端Listener服务
	go func() {
		// 1.获取TCP的Address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("服务启动失败：", err.Error())
			return
		}

		// 2.监听服务地址
		tcp, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("监听服务地址失败：", err.Error())
			return
		}

		// 3.启动server网络连接业务
		for true {
			// 等待客户端建立请求链接
			var conn *net.TCPConn
			conn, err = tcp.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP ERR：", err.Error())
				continue
			}

			// 建立连接成功
			fmt.Println("成功建立新的客户端连接 -> ", conn.RemoteAddr().String())

			// 我们这里暂时做一个最大512字节的回显服务
			go func() {
				// 不断的循环从客户端获取数据
				for {
					buf := make([]byte, 512)
					cnt, _ := conn.Read(buf)
					if cnt == 0 {
						continue
					}
					// 回显
					_, _ = conn.Write(buf[:cnt])
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("服务关闭")
}

func (s *Server) Server() {
	s.Start()

	// 阻塞主线程
	for true {
		time.Sleep(10 * time.Second)
	}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
