package znet

type Message struct {
	Id      uint32 // 消息ID
	DataLen uint32 // 消息长度
	Data    []byte // 消息内容
}

// NewMsgPackage 新建消息包
func NewMsgPackage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetDataLen(u uint32) {
	m.DataLen = u
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) SetMsgId(u uint32) {
	m.Id = u
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}
