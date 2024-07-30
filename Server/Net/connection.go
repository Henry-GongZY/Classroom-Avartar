package Net

import (
	"Server/NetInterface"
	"Server/utils"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type Connection struct {
	//链接所属服务器
	Server NetInterface.IServer
	//TCP链接
	Conn *net.TCPConn
	//链接ID
	ConnID uint32
	//是否关闭
	isClosed bool
	//停止标志管道
	ExitChan chan bool
	//读写协程消息通信
	msgChan chan []byte
	//报文解析器
	Task NetInterface.ITask
	//教师、学生姓名
	name string
	//教师,学生id
	userid uint32
	//由老师/学生使用的课程id
	lessonid uint32
	//由老师/学生使用的课程名
	lessonname string
	//判断是不是学生
	STU bool
	//教室
	ClassRoom NetInterface.IClassroom
}

func Int2String(length int) (string, error) {
	if length < 10 {
		return "Length:000" + strconv.Itoa(length), nil
	} else if length > 10 && length < 100 {
		return "Length:00" + strconv.Itoa(length), nil
	} else if length >= 100 && length < 1000 {
		return "Length:0" + strconv.Itoa(length), nil
	} else if length >= 1000 && length < 10000 {
		return "Length:" + strconv.Itoa(length), nil
	}
	return "", errors.New("Packet length error!")
}

func (c *Connection) GetServer() NetInterface.IServer {
	return c.Server
}

func NewConnection(Server NetInterface.IServer, conn *net.TCPConn, connID uint32, parser NetInterface.ITask) *Connection {
	c := &Connection{
		Server:   Server,
		Conn:     conn,
		ConnID:   connID,
		isClosed: false,
		ExitChan: make(chan bool, 1),
		Task:     parser,
		msgChan:  make(chan []byte),
		lessonid: 0,
	}

	//用管理器进行链接管理
	c.Server.GetConnMg().Add(c)

	return c
}

//从客户端接收消息的模块
func (c *Connection) StartReader() {
	fmt.Println("Reader Coroutine is running!")
	//业务结束逻辑
	defer fmt.Println("ConnID = ", c.ConnID, "Reader Coroutine exit, remote addr is ", c.RemoteAddr().String())
	defer c.Stop()

	for true {
		//创建拆包解包对象，以下过程将包解为3部分，并打包进msg
		dp := NewDataPack()
		//读取客户端的Msg Head 11字节
		headData := make([]byte, dp.GetHeadlen())
		_, err := io.ReadFull(c.GetTCPCpnnection(), headData)
		if err != nil {
			fmt.Println("Flushing connection buffer!")
			break
		}

		//解析头部
		str := string(headData)
		sli := strings.Split(str, ":")
		sli[1] = strings.TrimLeft(sli[1], "0")
		Len, err := strconv.ParseInt(sli[1], 10, 32)
		if err != nil {
			fmt.Println(err)
		}

		//根据头部数据继续读取
		databuf := make([]byte, Len)
		if _, err := io.ReadFull(c.GetTCPCpnnection(), databuf); err != nil {
			fmt.Println("Read message error:", err)
			break
		}

		//拆包
		msg, err := dp.Unpack(databuf)
		if err != nil {
			fmt.Println("Unpack error:", err)
			break
		}
		msg.SetMsgLen(Len)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		//报文处理,发送工作池或者发送处理句柄（判断消息队列是否启动完成）
		if utils.GlobalObject.WorkerPoolSize > 0 {
			c.Task.SendMsgToTaskQueue(&req)
		} else {
			go c.Task.StartHandler(&req)
		}
	}
}

//发送消息给客户端的模块
func (c *Connection) StartWriter() {
	fmt.Println("Writter Coroutine is running!")
	defer fmt.Println("ConnID = ", c.ConnID, "Writter Coroutine exit, remote addr is ", c.RemoteAddr().String())

	for {
		select {
		case data := <-c.msgChan:
			//准备头部
			k := len(data)
			HeadBuf, err := Int2String(k)
			if err != nil {
				fmt.Println(err)
			}
			pak := []byte(HeadBuf)
			//打包发送
			pak = append(pak, data...)
			if _, err := c.Conn.Write(pak); err != nil {
				fmt.Println("Send data error,", err)
				return
			}
		case <-c.ExitChan:
			return
		}
	}
}

//链接启动
func (c *Connection) Start() {
	//读数据业务启动
	go c.StartReader()
	//写数据业务启动
	go c.StartWriter()
}

//断开链接
func (c *Connection) Stop() {
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	err := c.Conn.Close()
	if err != nil {
		return
	}

	//学生关闭连接
	if c.STU && c.ClassRoom != nil {
		fmt.Println("A student quit!")
		err := c.GetTeaConn().SendMsg(9, 0, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		if err != nil {
			fmt.Println("Send msg error when student get away: ", err)
		}
		err = c.GetTeaConn().DeleteStudent(c.userid)
		if err != nil {
			fmt.Println("Student has been offline :", err)
		}
	}

	//老师关闭连接
	if !c.STU && c.ClassRoom != nil {
		fmt.Println("A teacher quit!")
		//for _, conn := range c.ClassRoom.GetStudents() {
		//	err := conn.SendMsg(0, 0, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		//	if err != nil {
		//		fmt.Println("Send msg error when teacher get away: ", err)
		//	}
		//}

		//同时老师把班级退出
		err := c.GetServer().GetConnMg().Delroom(c.ClassRoom.Getlessonid())
		if err != nil {
			fmt.Println("Error: no such a classroom!")
		}
	}

	c.Server.GetConnMg().Remove(c)

	c.ExitChan <- true

	close(c.ExitChan)
	close(c.msgChan)
}

//获取链接的绑定socket
func (c *Connection) GetTCPCpnnection() *net.TCPConn {
	return c.Conn
}

//获取远程TCP状态
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

//获取当前连接的ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

//发送数据进入管道
func (c *Connection) SendMsg(msgId uint8, msgId2 uint8, msgId3 string, data []byte, gest []float32) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg!")
	}

	//封包
	dp := NewDataPack()
	dataPack, err := dp.Pack(NewMsgPack(msgId, msgId2, msgId3, data, gest))
	if err != nil {
		fmt.Println("Pack error:", err)
		return errors.New("Pack message error.")
	}

	//发送给管道
	c.msgChan <- dataPack
	return nil
}

//老师开教室
func (c *Connection) AddClassroom(Lid uint32, Lname string) error {
	//没有教室，建立一个
	if c.ClassRoom == nil {
		c.ClassRoom = &Classroom{
			lessonid:   Lid,
			lessonname: Lname,
			teacher:    c,
			students:   make(map[uint32]NetInterface.IConnection),
			signtable:  make(map[uint32]bool),
		}
		return nil
	} else {
		return errors.New("Classroom already exists!")
	}

}

//老师离开教室
func (c *Connection) DeleteClassroom() error {
	return nil
}

//找到教室
func (c *Connection) GetClassroom() NetInterface.IClassroom {
	return c.ClassRoom
}

//教室加学生
func (c *Connection) AddStudent(connection NetInterface.IConnection, Stuid uint32) error {
	return c.ClassRoom.AddStudent(connection, Stuid)
}

//教室删除学生
func (c *Connection) DeleteStudent(Stuid uint32) error {
	return c.ClassRoom.DelStudent(Stuid)
}

func (c *Connection) Setuserid(id uint32) {
	c.userid = id
}

func (c *Connection) Getuserid() uint32 {
	return c.userid
}

func (c *Connection) Setname(name string) {
	c.name = name
}

func (c *Connection) Setlid(id uint32) {
	c.lessonid = id
}

func (c *Connection) Getlid() uint32 {
	return c.lessonid
}

func (c *Connection) SetClassroom(classroom NetInterface.IClassroom) {
	c.ClassRoom = classroom
}

func (c *Connection) GetTeaConn() NetInterface.IConnection {
	return c.ClassRoom.GetTeacher()
}

func (c *Connection) SetStudent(tf bool) {
	c.STU = tf
}

func (c *Connection) Isstudent() bool {
	return c.STU
}

func (c *Connection) Setlname(name string) {
	c.lessonname = name
}

func (c *Connection) Getlname() string {
	return c.lessonname
}
