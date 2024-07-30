package NetInterface

type IRequest interface {

	//获取链接
	GetConnection() IConnection

	//获取数据
	GetData() []byte

	//获取消息ID
	GetMsgID() uint8
	GetMsgID2() uint8
	GetMsgID3() string
	GetMesg() string

	//获取长度
	GetMsgLen() int64
	//获取姿态
	GetGesture() []float32
}
