package Net

import (
	"Server/NetInterface"
)

type Request struct {
	//链接
	conn NetInterface.IConnection
	//数据
	msg NetInterface.IMessage
}

//获取链接
func (r *Request) GetConnection() NetInterface.IConnection {
	return r.conn
}

//获取数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

//获取消息Id1
func (r *Request) GetMsgID() uint8 {
	return r.msg.GetMsgId()
}

//获取消息Id2
func (r *Request) GetMsgID2() uint8 {
	return r.msg.GetMsgId2()
}

//获取消息Id3
func (r *Request) GetMsgID3() string {
	return r.msg.GetMsgId3()
}

//获取消息长度
func (r *Request) GetMsgLen() int64 {
	return r.msg.GetMsgLen()
}

func (r *Request) GetGesture() []float32 {
	return r.msg.GetGesture()
}

func (r *Request) GetMesg() string {
	return r.msg.GetMesg()
}
