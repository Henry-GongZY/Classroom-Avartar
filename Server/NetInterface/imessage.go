package NetInterface

type IMessage interface {
	//获取消息id
	GetMsgId() uint8
	GetMsgId2() uint8
	GetMsgId3() string
	GetMesg() string

	//获取数据
	GetData() []byte

	//设置消息id
	SetMsgId(uint8)
	SetMsgId2(uint8)
	SetMsgId3(string)
	SetMesg(string)

	//设置消息数据
	SetData([]byte)

	//消息长度
	GetMsgLen() int64
	SetMsgLen(int64)

	//姿态更新
	GetGesture() []float32
	SetGesture([]float32)
}
