package iface

type IMessage interface {
	// GetDataLen 获取消息长度
	GetDataLen() uint32
	// SetDataLen 设置消息长度
	SetDataLen(uint32)

	// GetMsgId 获取消息ID
	GetMsgId() uint32
	// SetMsgId 设置消息ID
	SetMsgId(uint32)

	// GetData 获取消息内容
	GetData() []byte
	// SetData 设置消息内容
	SetData([]byte)
}
