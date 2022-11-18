package network

import "github.com/451008604/socketServerFrame/iface"

type Request struct {
	conn iface.IConnection // 已经和客户端建立好的连接
	msg  iface.IMessage    // 客户端请求的数据
}

// GetConnection 获取请求的连接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// GetMsgID 获取请求消息的ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
