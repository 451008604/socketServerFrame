package znet

import (
	"errors"
	"fmt"
	"net"
	"socketServerFrame/iface"
	"time"
)

// Server 定义Server服务类实现IServer接口
type Server struct {
	Name      string // 服务器名称
	IPVersion string // tcp4 or other
	IP        string // IP地址
	Port      int    // 服务端口
	connID    uint32 // 客户端连接自增ID
}

func ResToClient(conn *net.TCPConn, data []byte, cnt int) error {
	if _, err := conn.Write(data[:cnt]); err != nil {
		return errors.New("回复客户端失败")
	}
	return nil
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
			// 等待客户端建立请求连接
			var conn *net.TCPConn
			conn, err = tcp.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP ERR：", err.Error())
				continue
			}
			// 自增connID
			s.connID++
			// 建立连接成功
			fmt.Println("成功建立新的客户端连接 -> ", conn.RemoteAddr().String(), "connID - ", s.connID)

			// 建立新的连接并监听客户端请求的消息
			dealConn := NewConnection(conn, s.connID, ResToClient)
			go dealConn.Start()
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

func NewServer(name string) iface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
