package api

type PingReq struct {
	Msg       string `json:"msg"`
	TimeStamp int64  `json:"time_stamp"`
}
