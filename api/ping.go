package api

import (
	"github.com/451008604/socketServerFrame/iface"
	"github.com/451008604/socketServerFrame/network"
	pb "github.com/451008604/socketServerFrame/proto/bin"
	"time"
)

type PingRouter struct {
	network.BaseRouter
}

func (p *PingRouter) Handler(req iface.IRequest) {
	t := time.Now().UnixMicro()
	pingReq := &pb.Ping{}
	UnmarshalProtoData(req.GetData(), pingReq)
	pingReq.TimeStamp = t - pingReq.GetTimeStamp()

	req.GetConnection().SendMsg(uint32(pb.MessageID_PING), MarshalProtoData(pingReq))
}
