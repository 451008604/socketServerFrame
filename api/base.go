package api

import (
	"encoding/json"

	"github.com/451008604/socketServerFrame/logs"
	"google.golang.org/protobuf/proto"
)

func ProtocolToByte(str proto.Message) []byte {
	marshal, err := proto.Marshal(str)
	if err != nil {
		logs.PrintLogErr(err)
		return []byte{}
	}
	return marshal
}

func ByteToProtocol(byte []byte, target proto.Message) {
	err := json.Unmarshal(byte, target)
	if err != nil {
		logs.PrintLogErr(err)
		return
	}
}
