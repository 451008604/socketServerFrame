package iface

// IDataPack 封包拆包，通过固定的包头获取消息数据，解决TCP粘包问题
type IDataPack interface {
	GetHeadLen() uint32       // 获取包头长度
	Pack(msg IMessage) []byte // 消息拆包
	Unpack([]byte) IMessage   // 消息封包
}
