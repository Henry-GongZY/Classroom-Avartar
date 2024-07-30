package Handles

import (
	"Server/Net"
	"Server/NetInterface"
	"fmt"
	"strconv"
)

type LoginHandler struct {
	Net.Handler
}

func (this *LoginHandler) Handle(request NetInterface.IRequest) {
	SorTid, _ := strconv.Atoi(request.GetMsgID3())
	Password := string(request.GetData())
	switch request.GetMsgID2() {
	//学生，去查学生表
	case 0:
		{
			cursor := request.GetConnection().GetServer().GetDbCursor()
			if found, name := cursor.StuLogin(uint32(SorTid), Password); found {
				//登录成功
				err := request.GetConnection().SendMsg(0, 8, request.GetMsgID3(), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				request.GetConnection().Setuserid(uint32(SorTid))
				request.GetConnection().Setname(name)
				request.GetConnection().SetStudent(true)
				if err != nil {
					fmt.Println("Student login success! Sendback error: ", err)
					return
				}
			} else {
				err := request.GetConnection().SendMsg(0, 9, request.GetMsgID3(), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Student login failed! Sendback error: ", err)
					return
				}
			}
			break
		}
	//教师，去查教师表
	case 1:
		{
			cursor := request.GetConnection().GetServer().GetDbCursor()
			if found, name := cursor.TeaLogin(uint32(SorTid), Password); found {
				err := request.GetConnection().SendMsg(0, 8, request.GetMsgID3(), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				request.GetConnection().Setuserid(uint32(SorTid))
				request.GetConnection().Setname(name)
				request.GetConnection().SetStudent(false)
				if err != nil {
					fmt.Println("Teacher login success! Sendback error: ", err)
					return
				}
			} else {
				err := request.GetConnection().SendMsg(0, 9, request.GetMsgID3(), []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Teacher login failed! Sendback error: ", err)
					return
				}
			}
			break
		}
	default:
		fmt.Println("Error in Id2!")
	}

}
