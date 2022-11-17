package api

import (
	"socketServerFrame/iface"
	pb "socketServerFrame/proto/bin"
	"socketServerFrame/znet"
	"time"
)

type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handler(req iface.IRequest) {
	t := time.Now().UnixMicro()
	pingReq := &pb.Ping{}
	UnmarshalProtoData(req.GetData(), pingReq)
	pingReq.TimeStamp = t - pingReq.GetTimeStamp()

	req.GetConnection().SendMsg(uint32(pb.MessageID_PING), MarshalProtoData(pingReq))
}
