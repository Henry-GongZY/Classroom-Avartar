package NetInterface

type ITask interface {

	//调度相应的执行程序
	StartHandler(request IRequest)
	//添加处理逻辑
	AddHandler(msgId uint8, handler IHandler)
	//启动工作池
	StartWorkerPool()
	//向任务队列发送消息
	SendMsgToTaskQueue(IRequest)
}
