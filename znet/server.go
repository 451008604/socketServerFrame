package znet

import (
	"fmt"
	"net"
	"socketServerFrame/config"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
	"time"
)

// Server 定义Server服务类实现IServer接口
type Server struct {
	Name       string            // 服务器名称
	IPVersion  string            // tcp4 or other
	IP         string            // IP地址
	Port       string            // 服务端口
	msgHandler iface.IMsgHandler // 当前Server的消息管理模块，用来绑定MsgId和对应的处理函数
	connID     uint32            // 客户端连接自增ID
}

func NewServer() iface.IServer {
	s := &Server{
		Name:       config.GetGlobalObject().Name,
		IPVersion:  "tcp4",
		IP:         config.GetGlobalObject().Host,
		Port:       config.GetGlobalObject().TcpPort,
		msgHandler: NewMsgHandler(),
		connID:     0,
	}
	return s
}

func (s *Server) Start() {
	// 开启一个go去做服务端Listener服务
	go func() {
		// 启动工作池等待接收请求数据
		s.msgHandler.StartWorkerPool()

		// 1.获取TCP的Address
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%s", s.IP, s.Port))
		if logs.PrintLogErrToConsole(err, "服务启动失败：") {
			return
		}

		// 2.监听服务地址
		tcp, err := net.ListenTCP(s.IPVersion, addr)
		if logs.PrintLogErrToConsole(err, "监听服务地址失败：") {
			return
		}

		// 3.启动server网络连接业务
		for true {
			// 等待客户端建立请求连接
			var conn *net.TCPConn
			conn, err = tcp.AcceptTCP()
			if logs.PrintLogErrToConsole(err, "AcceptTCP ERR：") {
				continue
			}
			// 自增connID
			s.connID++
			// 建立连接成功
			logs.PrintLogInfoToConsole(fmt.Sprintf("成功建立新的客户端连接 -> %v connID - %v", conn.RemoteAddr().String(), s.connID))

			// 建立新的连接并监听客户端请求的消息
			dealConn := NewConnection(conn, s.connID, s.msgHandler)
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

func (s *Server) AddRouter(msgId uint32, router iface.IRouter) {
	s.msgHandler.AddRouter(msgId, router)
}
