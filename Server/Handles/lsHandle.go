package Handles

import (
	"Server/Data"
	"Server/Net"
	"Server/NetInterface"
	"fmt"
	"strconv"
)

type LsHandler struct {
	Net.Handler
}

func Format(data []Data.Lesson) []byte {
	str := ""
	for i := range data {
		str += data[i].Lessonname
		str += "\n"
	}
	return []byte(str)
}

func (this *LsHandler) Handle(request NetInterface.IRequest) {
	db := request.GetConnection().GetServer().GetDbCursor()
	id2 := request.GetMsgID2()
	switch id2 {
	case 0:
		{
			//查找学生正在进行的课程
			SorTid, _ := strconv.Atoi(request.GetMsgID3())
			if lessons, t := db.StuClasses(uint32(SorTid)); t {
				err := request.GetConnection().SendMsg(5, 2, request.GetMsgID3(), Format(lessons), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println(err)
				}
			}
			break
		}
	case 1:
		{
			//查找老师正在进行的课程
			SorTid, _ := strconv.Atoi(request.GetMsgID3())
			if lessons, t := db.TeaClasses(uint32(SorTid)); t {
				err := request.GetConnection().SendMsg(5, 2, request.GetMsgID3(), Format(lessons), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println(err)
				}
			}
			break
		}
	default:
		fmt.Println("Error in LessonHandle: shouldn't be here!")
	}
}
