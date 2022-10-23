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
	Name        string                             // 服务器名称
	IPVersion   string                             // tcp4 or other
	IP          string                             // IP地址
	Port        string                             // 服务端口
	msgHandler  iface.IMsgHandler                  // 当前Server的消息管理模块，用来绑定MsgId和对应的处理函数
	connMgr     iface.IConnManager                 // 当前Server的连接管理器
	OnConnStart func(connection iface.IConnection) // 该Server连接创建时的Hook函数
	OnConnStop  func(connection iface.IConnection) // 该Server连接断开时的Hook函数
	connID      uint32                             // 客户端连接自增ID
}

func NewServer() iface.IServer {
	s := &Server{
		Name:        config.GetGlobalObject().Name,
		IPVersion:   "tcp4",
		IP:          config.GetGlobalObject().Host,
		Port:        config.GetGlobalObject().TcpPort,
		msgHandler:  NewMsgHandler(),
		connMgr:     NewConnManager(),
		OnConnStart: nil,
		OnConnStop:  nil,
		connID:      0,
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

			// 连接数量超过限制后，关闭新建立的连接
			if s.connMgr.Len() >= config.GetGlobalObject().MaxConn {
				_ = conn.Close()
				continue
			}

			// 自增connID
			s.connID = uint32(s.GetConnMgr().Len() + 1)
			// 建立连接成功
			logs.PrintLogInfoToConsole(fmt.Sprintf("成功建立新的客户端连接 -> %v connID - %v", conn.RemoteAddr().String(), s.connID))

			// 建立新的连接并监听客户端请求的消息
			dealConn := NewConnection(s, conn, s.connID, s.msgHandler)
			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("服务关闭")

	s.connMgr.ClearConn()
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

func (s *Server) GetConnMgr() iface.IConnManager {
	return s.connMgr
}

// SetOnConnStart Server连接创建时的Hook函数
func (s *Server) SetOnConnStart(f func(conn iface.IConnection)) {
	s.OnConnStart = f
}

// CallbackOnConnStart 调用Server连接时的Hook函数
func (s *Server) CallbackOnConnStart(conn iface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

// SetOnConnStop Server连接断开时的Hook函数
func (s *Server) SetOnConnStop(f func(conn iface.IConnection)) {
	s.OnConnStop = f
}

// CallbackOnConnStop 调用Server连接断开时的Hook函数
func (s *Server) CallbackOnConnStop(conn iface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}
