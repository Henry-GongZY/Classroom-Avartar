package Net

import (
	"Server/NetInterface"
)

type Handler struct{}

//处理连接业务之前的方法
func (r *Handler) PreHandle(request NetInterface.IRequest) {}

//处理连接业务的主方法
func (r *Handler) Handle(request NetInterface.IRequest) {}

//处理连接任务之后的方法
func (r *Handler) PostHandle(request NetInterface.IRequest) {}
