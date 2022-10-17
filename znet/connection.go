package znet

import (
	"io"
	"net"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
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
		dp := NewDataPack()

		// 获取客户端的消息头信息
		headData := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.GetTCPConnection(), headData)
		if logs.PrintToConsoleErr(err) {
			c.ExitBuffChan <- true
		}
		// 通过消息头获取dataLen和Id
		msgData := dp.Unpack(headData)
		if msgData == nil {
			c.ExitBuffChan <- true
		}
		// 通过消息头获取消息body
		if msgData.GetDataLen() > 0 {
			msgData.SetData(make([]byte, msgData.GetDataLen()))
			_, err = io.ReadFull(c.GetTCPConnection(), msgData.GetData())
			if logs.PrintToConsoleErr(err) {
				c.ExitBuffChan <- true
				continue
			}
		}

		// 封装请求和请求数据
		req := &Request{conn: c, msg: msgData}
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
	if c.isClosed {
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

// SendMsg 发送消息给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) {
	if c.isClosed {
		logs.PrintToConsoleInfo("连接已关闭导致消息发送失败")
		return
	}

	// 新建数据传输包
	dp := NewDataPack()
	// 将消息数据封包
	msg := dp.Pack(NewMsgPackage(msgId, data))
	if msg == nil {
		return
	}
	// 写入传输通道发送给客户端
	_, err := c.Conn.Write(msg)
	if logs.PrintToConsoleErr(err) {
		c.ExitBuffChan <- true
	}
}
