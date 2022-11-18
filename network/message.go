package network

type Message struct {
	id      uint32 // 消息ID
	dataLen uint32 // 消息长度
	data    []byte // 消息内容
}

// NewMsgPackage 新建消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		id:      id,
		dataLen: uint32(len(data)),
		data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.dataLen
}

func (m *Message) SetDataLen(u uint32) {
	m.dataLen = u
}

func (m *Message) GetMsgId() uint32 {
	return m.id
}

func (m *Message) SetMsgId(u uint32) {
	m.id = u
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) SetData(bytes []byte) {
	m.data = bytes
}
