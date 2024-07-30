package Handles

import (
	"Server/Net"
	"Server/NetInterface"
	"fmt"
	"strconv"
	"time"
)

//重写句柄
type InfoHandler struct {
	Net.Handler
}

func (this *InfoHandler) Handle(request NetInterface.IRequest) {
	SorTid, _ := strconv.Atoi(request.GetMsgID3())
	fmt.Println(request.GetMsgID2())

	switch request.GetMsgID2() {
	//学生查找课程
	case 0:
		{
			lid, lesson, found := request.GetConnection().GetServer().GetDbCursor().StuClass(uint32(SorTid))
			time.Sleep(time.Millisecond * 100)
			_, f := request.GetConnection().GetServer().GetConnMg().Getroom(lid)
			//有课程，返回课程
			if !f && found {
				err := request.GetConnection().SendMsg(2, 8, request.GetMsgID3(), []byte(lesson), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
				//有课程且已经开课，返回房间
			} else if f && found {
				request.GetConnection().Setlid(lid)
				err := request.GetConnection().SendMsg(2, 10, request.GetMsgID3(), []byte(lesson), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
				//没有课程，返回9错误码
			} else {
				err := request.GetConnection().SendMsg(2, 9, request.GetMsgID3(), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
			}
			break
		}
	//教师查找课程
	case 1:
		{
			lid, lesson, found := request.GetConnection().GetServer().GetDbCursor().TeaClass(uint32(SorTid))
			time.Sleep(time.Millisecond * 100)
			//有课程，返回课程
			if found {
				//设置课程id和课程名称
				request.GetConnection().Setlid(lid)
				request.GetConnection().Setlname(lesson)
				err := request.GetConnection().SendMsg(2, 8, "", []byte(lesson), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
				//没有课程，返回9错误码
			} else {
				err := request.GetConnection().SendMsg(2, 9, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
			}
			break
		}
	//教师注册班级
	case 2:
		{
			//找到教室列表并注册教室
			err := request.GetConnection().AddClassroom(request.GetConnection().Getlid(), request.GetConnection().Getlname())
			_ = request.GetConnection().GetServer().GetConnMg().Addroom(request.GetConnection().Getlid(), request.GetConnection(), request.GetConnection().GetClassroom())
			if err != nil {
				//教室已经存在，注册失败
				fmt.Println("Error: ", err)
				err := request.GetConnection().SendMsg(2, 9, "", []byte("Classroom"), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
			} else {
				//教室不存在，注册成功
				err := request.GetConnection().SendMsg(2, 8, "", []byte("Classroom"), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
			}
			break
		}
	//学生加入班级
	case 3:
		{
			//课程
			if request.GetConnection().Getlid() == 0 {
				err := request.GetConnection().SendMsg(2, 4, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msg error: ", err)
				}
			} else {
				//学生加入成功
				classroom, _ := request.GetConnection().GetServer().GetConnMg().Getroom(request.GetConnection().Getlid())
				request.GetConnection().SetClassroom(classroom)
				err := request.GetConnection().GetClassroom().AddStudent(request.GetConnection(), request.GetConnection().Getuserid())
				if err == nil {
					err := request.GetConnection().SendMsg(2, 6, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
					if err != nil {
						fmt.Println("Send msg error back to student: ", err)
					}
					err = nil
					err = request.GetConnection().GetTeaConn().SendMsg(1, 1, strconv.Itoa(int(request.GetConnection().Getuserid())), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
					if err != nil {
						fmt.Println("Send msg error to teacher: ", err)
					}
				} else {
					//学生在班级里，加入失败
					fmt.Println(err)
					err = request.GetConnection().SendMsg(2, 5, "", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
					if err != nil {
						fmt.Println("Send msg error: ", err)
					}
				}
			}
			break
		}
	default:
		{
			fmt.Println("Error in Id2!")
			return
		}
	}

}
