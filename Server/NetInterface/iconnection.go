package NetInterface

import "net"

//连接和释放的抽象层
type IConnection interface {
	//链接启动
	Start()

	//断开链接
	Stop()

	//获取链接的绑定socket
	GetTCPCpnnection() *net.TCPConn

	//获取远程TCP状态
	RemoteAddr() net.Addr

	//获取当前连接的ID
	GetConnID() uint32

	//发送数据给远程的客户端
	SendMsg(msgId uint8, msgId2 uint8, msgId3 string, data []byte, gest []float32) error

	//获取链接所属服务器
	GetServer() IServer

	//给老师绑定一个教师
	AddClassroom(uint32, string) error

	//推出教室
	DeleteClassroom() error

	//得到学生、老师id
	Setuserid(uint32)

	Getuserid() uint32

	Setname(string)

	DeleteStudent(uint32) error

	AddStudent(IConnection, uint32) error

	Setlid(uint32)

	Getlid() uint32

	SetClassroom(IClassroom)

	GetTeaConn() IConnection

	GetClassroom() IClassroom

	SetStudent(bool)

	Isstudent() bool

	Setlname(name string)

	Getlname() string
}

//链接业务处理方法
type HandleFunc func(*net.TCPConn, []byte, int) error
