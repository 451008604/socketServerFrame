package znet

import (
	"net"
	"socketServerFrame/iface"
)

type Connection struct {
	// 当前连接的SocketTCP套接字
	Conn *net.TCPConn
	// 当前连接的ID（SessionID）
	ConnID uint32
	// 当前连接是否已关闭
	isClosed bool
	// 该连接的处理方法router
	Router iface.IRouter
	// 通知该连接已经退出的channel
	ExitBuffChan chan bool
}

// NewConnection 新建连接
func NewConnection(conn *net.TCPConn, connID uint32, router iface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// StartReader 处理conn接收的客户端数据
func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		// 缓冲区数据写入到buf中
		_, err := c.Conn.Read(buf)
		if err != nil {
			c.ExitBuffChan <- true
			continue
		}

		// 封装请求和请求数据
		req := &Request{conn: c, data: buf}
		// 使用goroutine处理请求数据
		go func(request iface.IRequest) {
			c.Router.PreHandler(request)
			c.Router.Handler(request)
			c.Router.AfterHandler(request)
		}(req)
	}
}

// Start 启动连接
func (c *Connection) Start() {
	// 开启监听收到该连接请求数据后的处理
	go c.StartReader()

	for {
		select {
		// 在收到退出消息时释放进程
		case <-c.ExitBuffChan:
			return
		}
	}
}

// Stop 停止连接
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	_ = c.Conn.Close()
	// 通知关闭该连接的监听
	c.ExitBuffChan <- true
	// 关闭该连接管道
	close(c.ExitBuffChan)
}

// GetTCPConnection 从当前连接获取原始的Socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}
