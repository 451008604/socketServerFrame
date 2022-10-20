package znet

import (
	"io"
	"net"
	"socketServerFrame/config"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
)

type Connection struct {
	TcpServer    iface.IServer     // 当前Conn所属的Server
	Conn         *net.TCPConn      // 当前连接的SocketTCP套接字
	ConnID       uint32            // 当前连接的ID（SessionID）
	isClosed     bool              // 当前连接是否已关闭
	MsgHandler   iface.IMsgHandler // 消息管理MsgId和对应处理函数的消息管理模块
	ExitBuffChan chan bool         // 通知该连接已经退出的channel
	msgChan      chan []byte       // 用于读、写两个goroutine之间的消息通信
}

// NewConnection 新建连接
func NewConnection(server iface.IServer, conn *net.TCPConn, connID uint32, msgHandler iface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:    server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}

	// 将新建的连接添加到所属Server的连接管理器内
	c.TcpServer.GetConnMgr().Add(c)
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
		if logs.PrintLogErrToConsole(err) {
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
			if logs.PrintLogErrToConsole(err) {
				c.ExitBuffChan <- true
				continue
			}
		}

		// 封装请求和请求数据
		req := &Request{conn: c, msg: msgData}
		if config.GetGlobalObject().WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			go c.MsgHandler.DoMsgHandler(req)
		}
	}
}

// StartWriter 写消息goroutine，用户将数据发送给客户端
func (c *Connection) StartWriter() {
	for true {
		select {
		case data := <-c.msgChan:
			_, err := c.Conn.Write(data)
			if logs.PrintLogErrToConsole(err) {
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

// Start 启动连接
func (c *Connection) Start() {
	// 开启用于读的goroutine
	go c.StartReader()
	// 开启用于写的goroutine
	go c.StartWriter()

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
	// 将连接从连接管理器中删除
	c.TcpServer.GetConnMgr().Remove(c)

	// 关闭该连接管道
	close(c.ExitBuffChan)
	close(c.msgChan)
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
		logs.PrintLogInfoToConsole("连接已关闭导致消息发送失败")
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
	c.msgChan <- msg
}
