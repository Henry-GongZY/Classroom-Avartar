package Net

import (
	"Server/NetInterface"
	"Server/utils"
	"fmt"
	"strconv"
)

type Task struct {
	//不同ID对应的处理方法
	handles map[uint8]NetInterface.IHandler
	//消息队列
	TaskQueue []chan NetInterface.IRequest
	//Worker数量
	WorkerPoolSize uint32
}

//构造函数
func NewParser() *Task {
	return &Task{
		handles:        make(map[uint8]NetInterface.IHandler),
		TaskQueue:      make([]chan NetInterface.IRequest, utils.GlobalObject.WorkerPoolSize),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
	}
}

//根据发送来的Id启动对应的应对方法
func (p *Task) StartHandler(request NetInterface.IRequest) {
	handler, succ := p.handles[request.GetMsgID()]
	if !succ {
		fmt.Println("msgId = ", request.GetMsgID(), "is not found")
	}

	//handle运行
	handler.PreHandle(request)
	go handler.Handle(request)
	handler.PostHandle(request)
}

//加入处理方法
func (p *Task) AddHandler(msgId uint8, handler NetInterface.IHandler) {
	if _, succ := p.handles[msgId]; succ {
		panic("msgId = " + strconv.Itoa(int(msgId)))
	}
	p.handles[msgId] = handler
	fmt.Println("msgId = ", msgId, " added seccessfully.")
}

//启动Worker工作池
func (p *Task) StartWorkerPool() {
	for i := 0; i < int(p.WorkerPoolSize); i++ {
		p.TaskQueue[i] = make(chan NetInterface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		go p.StartOneWorker(i, p.TaskQueue[i])
	}
}

//等待消息，处理业务
func (p *Task) StartOneWorker(workerid int, taskQueue chan NetInterface.IRequest) {
	fmt.Println("WorkerID = ", workerid, " is started...")

	for {
		select {
		case request := <-taskQueue:
			p.StartHandler(request)
		}
	}
}

//将消息交给TaskQueue
func (p *Task) SendMsgToTaskQueue(request NetInterface.IRequest) {
	//使用最基础的轮询分配策略
	workerid := request.GetConnection().GetConnID() % p.WorkerPoolSize
	fmt.Println("Add ConnId = ", request.GetConnection().GetConnID(),
		" request msgId = ", request.GetMsgID(), " to workerId = ", workerid)

	//将连接分配给任务队列处理
	p.TaskQueue[workerid] <- request
}
