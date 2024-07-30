package Net

import (
	"Server/NetInterface"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type ConnManager struct {
	//链接管理,对应链接Id和链接本身
	connections map[uint32]NetInterface.IConnection
	//读写保护
	connLock sync.RWMutex
	//课程名单
	lelist map[uint32]NetInterface.IClassroom
	//读写保护
	classLock sync.RWMutex
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]NetInterface.IConnection),
		lelist:      make(map[uint32]NetInterface.IClassroom),
	}
}

func (cm *ConnManager) Add(conn NetInterface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//添加链接
	cm.connections[conn.GetConnID()] = conn
	fmt.Println("Connection added to ConnManager, Count = ", cm.Count(),
		" ConnId = ", conn.GetConnID())
}

func (cm *ConnManager) Remove(conn NetInterface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//删除链接
	delete(cm.connections, conn.GetConnID())
	fmt.Println("Connection removed , Count = ", cm.Count())
}

func (cm *ConnManager) GetConn(connID uint32) (NetInterface.IConnection, error) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	//搜索链接
	if conn, found := cm.connections[connID]; found {
		return conn, nil
	} else {
		return nil, errors.New("Connection not found!")
	}
}

func (cm *ConnManager) Count() int {
	return len(cm.connections)
}

func (cm *ConnManager) Clear() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connID, conn := range cm.connections {
		//停止链接
		conn.Stop()
		//释放资源
		delete(cm.connections, connID)
	}

	fmt.Println("All connections cleared!")
}

//添加房间
func (cm *ConnManager) Addroom(Lessonid uint32, conn NetInterface.IConnection, classroom NetInterface.IClassroom) error {
	//加锁
	cm.classLock.Lock()
	defer cm.classLock.Unlock()

	//还没注册
	if cm.lelist[Lessonid] == nil {
		cm.lelist[Lessonid] = classroom

		go func() {
			classroom.GenStulist(conn.GetServer().GetDbCursor().Students(Lessonid))
		}()

		return nil
		//完成注册
	} else {
		return errors.New("Registered!")
	}
}

//得到房间
func (cm *ConnManager) Getroom(Lessonid uint32) (NetInterface.IClassroom, bool) {
	cm.classLock.Lock()
	defer cm.classLock.Unlock()

	if cm.lelist[Lessonid] == nil {
		return nil, false
	} else {
		return cm.lelist[Lessonid], true
	}
}

//删除房间
func (cm *ConnManager) Delroom(Lessonid uint32) error {
	cm.classLock.Lock()
	defer cm.classLock.Unlock()

	if cm.lelist[Lessonid] != nil {
		l := cm.lelist[Lessonid]
		go func() {
			l.Signin()
		}()
		delete(cm.lelist, Lessonid)
		fmt.Println("Room " + strconv.Itoa(int(Lessonid)) + " is deleted!")
		return nil
	} else {
		return errors.New("No such room need to be deleted!")
	}
}
