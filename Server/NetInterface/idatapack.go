package NetInterface

type IDatapack interface {
	//获取封包长度
	GetHeadlen() uint32
	//封包
	Pack(msg IMessage) ([]byte, error)
	//拆包
	Unpack([]byte) (IMessage, error)
}
