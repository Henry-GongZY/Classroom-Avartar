package NetInterface

type IConnManager interface {
	//增加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//修改链接
	GetConn(connID uint32) (IConnection, error)
	//得到链接数目
	Count() int
	//清除链接
	Clear()
	//注册课程
	Addroom(uint32, IConnection, IClassroom) error
	//查找房间
	Getroom(uint32) (IClassroom, bool)
	//删除房间
	Delroom(uint32) error
}
