package znet

import (
	"bytes"
	"encoding/binary"
	"socketServerFrame/config"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
)

type DataPack struct{}

func (d *DataPack) GetHeadLen() uint32 {
	// id uint32(4字节) + dataLen uint32(4字节)
	return 8
}

// NewDataPack 新数据包
func NewDataPack() *DataPack {
	return &DataPack{}
}

// Pack 封包
func (d *DataPack) Pack(msg iface.IMessage) []byte {
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if logs.PrintToConsoleErr(binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen())) {
		return nil
	}
	// 写msgId
	if logs.PrintToConsoleErr(binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId())) {
		return nil
	}
	// 写data数据
	if logs.PrintToConsoleErr(binary.Write(dataBuff, binary.LittleEndian, msg.GetData())) {
		return nil
	}
	return dataBuff.Bytes()
}

// Unpack 拆包(只获取到包头Id,dataLen)
func (d *DataPack) Unpack(binaryData []byte) iface.IMessage {
	dataBuff := bytes.NewReader(binaryData)
	msgData := &Message{}

	// 读dataLen
	if logs.PrintToConsoleErr(binary.Read(dataBuff, binary.LittleEndian, &msgData.dataLen)) {
		return nil
	}
	// 读msgId
	if logs.PrintToConsoleErr(binary.Read(dataBuff, binary.LittleEndian, &msgData.id)) {
		return nil
	}
	// 检查数据长度是否超出限制
	if config.GetGlobalObject().MaxPackSize > 0 && msgData.GetDataLen() > config.GetGlobalObject().MaxPackSize {
		logs.PrintToConsoleInfo("接收数据长度超限")
		return nil
	}
	return msgData
}
