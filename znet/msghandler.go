package znet

import (
	"errors"
	"fmt"
	"socketServerFrame/iface"
	"socketServerFrame/logs"
)

type MsgHandler struct {
	Apis map[uint32]iface.IRouter // 存放每个MsgId所对应处理方法的map属性
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{Apis: make(map[uint32]iface.IRouter)}
}

// DoMsgHandler 执行路由绑定的处理函数
func (m *MsgHandler) DoMsgHandler(request iface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		logs.PrintLogInfoToConsole(fmt.Sprintf("api msgID %v is not fund", request.GetMsgID()))
		return
	}

	// 对应的逻辑处理方法
	handler.PreHandler(request)
	handler.Handler(request)
	handler.AfterHandler(request)
}

// AddRouter 添加路由，绑定处理函数
func (m *MsgHandler) AddRouter(msgId uint32, router iface.IRouter) {
	if _, ok := m.Apis[msgId]; ok {
		logs.PrintLogPanicToConsole(errors.New("消息ID重复绑定Handler"))
	}
	m.Apis[msgId] = router
}
