package network

import (
	"errors"
	"fmt"

	"github.com/451008604/socketServerFrame/config"
	"github.com/451008604/socketServerFrame/iface"
	"github.com/451008604/socketServerFrame/logs"
)

type MsgHandler struct {
	Apis           map[uint32]iface.IRouter // 存放每个MsgId所对应处理方法的map属性
	WorkerPoolSize uint32                   // 业务工作Work池的数量
	TaskQueue      []chan iface.IRequest    // Worker负责取任务的消息队列
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]iface.IRouter),
		WorkerPoolSize: config.GetGlobalObject().WorkerPoolSize,
		TaskQueue:      make([]chan iface.IRequest, config.GetGlobalObject().WorkerPoolSize),
	}
}

// DoMsgHandler 执行路由绑定的处理函数
func (m *MsgHandler) DoMsgHandler(request iface.IRequest) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		logs.PrintLogInfo(fmt.Sprintf("api msgID %v is not fund", request.GetMsgID()))
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
		logs.PrintLogPanic(errors.New("消息ID重复绑定Handler"))
	}
	m.Apis[msgId] = router
}

// StartWorkerPool 启动工作池
func (m *MsgHandler) StartWorkerPool() {
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		m.TaskQueue[i] = make(chan iface.IRequest, config.GetGlobalObject().WorkerTaskMaxLen)

		go m.StartOneWorker(m.TaskQueue[i])
	}
}

// SendMsgToTaskQueue 将消息发送到任务队列
func (m *MsgHandler) SendMsgToTaskQueue(request iface.IRequest) {
	// 根据connID平均分配至对应worker
	workerID := request.GetConnection().GetConnID() % m.WorkerPoolSize
	// 将请求推入worker协程
	m.TaskQueue[workerID] <- request
}

// StartOneWorker 启动一个工作协程等待处理接收的请求
func (m *MsgHandler) StartOneWorker(taskQueue chan iface.IRequest) {
	for request := range taskQueue {
		m.DoMsgHandler(request)
	}
}
