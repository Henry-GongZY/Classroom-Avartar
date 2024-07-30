package Net

import (
	"Server/NetInterface"
	"Server/utils"
	"errors"
	"sync"
)

type Classroom struct {
	//课程id
	lessonid uint32
	//课程名称
	lessonname string
	//读写锁
	roomlock sync.RWMutex
	//老师连接
	teacher NetInterface.IConnection
	//学生连接
	students map[uint32]NetInterface.IConnection
	//签名表
	signtable map[uint32]bool
}

//找老师
func (c *Classroom) GetTeacher() NetInterface.IConnection {
	return c.teacher
}

//找学生
func (c *Classroom) GetStudent(Stuid uint32) NetInterface.IConnection {
	if c.students[Stuid] == nil {
		//没有学生
		return nil
	} else {
		//有学生
		return c.students[Stuid]
	}
}

//加学生
func (c *Classroom) AddStudent(connection NetInterface.IConnection, Stuid uint32) error {
	//加减锁
	c.roomlock.Lock()
	defer c.roomlock.Unlock()

	if c.GetStudent(Stuid) != nil {
		return errors.New("Student's already in the classroom!")
	} else {
		c.students[Stuid] = connection
		//签到
		c.signtable[Stuid] = true
		return nil
	}
}

//删学生
func (c *Classroom) DelStudent(Stuid uint32) error {
	//加减锁
	c.roomlock.Lock()
	defer c.roomlock.Unlock()

	if c.students[Stuid] != nil {
		delete(c.students, Stuid)
		return nil
	} else {
		return errors.New("Student has already left!")
	}
}

//得到学生表
func (c *Classroom) GetStudents() map[uint32]NetInterface.IConnection {
	return c.students
}

//设置课程Id
func (c *Classroom) SetLessonid(lessonid uint32) {
	c.lessonid = lessonid
}

//设置课程名称
func (c *Classroom) SetLessonname(lessonname string) {
	c.lessonname = lessonname
}

//获取课程Id
func (c *Classroom) Getlessonid() uint32 {
	return c.lessonid
}

//获取课程名称
func (c *Classroom) Getlessonname() string {
	return c.lessonname
}

//生成学生名单
func (c *Classroom) GenStulist(stus []uint32) {
	for i := 0; i <= len(stus)-1; i++ {
		c.signtable[stus[i]] = false
	}
}

//完成签到
func (c *Classroom) Signin() {
	s := utils.SigninObj{}
	s.Gensignin(c.signtable, c.lessonname)
}
