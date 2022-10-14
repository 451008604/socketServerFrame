package znet

import "socketServerFrame/iface"

type Request struct {
	conn iface.IConnection // 已经和客户端建立好的连接
	data []byte            // 客户端请求的数据
}

// GetConnection 获取请求的连接信息
func (r *Request) GetConnection() iface.IConnection {
	return r.conn
}

// GetData 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.data
}
