package Handles

import (
	"Server/Data"
	"Server/ManagementInterface"
	"Server/Net"
	"Server/NetInterface"
	"fmt"
	"strconv"
	"strings"
)

type FileUploadHandler struct {
	Net.Handler
	fh       ManagementInterface.IFiletask
	Handleon bool
}

func (this *FileUploadHandler) PreHandle(request NetInterface.IRequest) {
	if !this.Handleon {
		this.fh = request.GetConnection().GetServer().GetFilehandler()
		this.Handleon = true
	}
}

func (this *FileUploadHandler) Handle(request NetInterface.IRequest) {
	//学生或老师的id
	SorTid := request.GetMsgID3()
	data := request.GetData()
	//文件名+课程名
	lessonname := request.GetMesg()
	command := int(request.GetMsgID2())

	switch command {
	case 1:
	case 2:
	case 3:
		this.fh.GetChannel() <- Data.FileCommand{data, lessonname + "/" + SorTid, command}
		break
	case 4:
		//查找文件
		lessonname = strings.Split(lessonname, "/")[1]
		//学生文件
		files := ManagementInterface.GetDoclist(lessonname, SorTid, request.GetConnection().Isstudent())
		filestr := ""
		for i := range files {
			filestr += (files[i] + "\n")
		}
		//教师文件
		teafiles := ManagementInterface.GetTeaDoclist(lessonname, SorTid, request.GetConnection().Isstudent())
		teafilestr := ""
		for i := range teafiles {
			teafilestr += (teafiles[i] + "\n")
		}

		err := request.GetConnection().SendMsg(4, 6, strconv.Itoa(int(request.GetConnection().Getuserid())), []byte(filestr), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		if err != nil {
			fmt.Println(err)
		}

		err = nil
		err = request.GetConnection().SendMsg(4, 7, strconv.Itoa(int(request.GetConnection().Getuserid())), []byte(teafilestr), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
		if err != nil {
			fmt.Println(err)
		}

		break
	case 5:
		//删除文件
		strsect := strings.Split(lessonname, "/")
		filename := strsect[0]
		lessonname = strsect[1]
		ManagementInterface.DelDoclist(lessonname, SorTid, request.GetConnection().Isstudent(), filename)
		break
	default:
		fmt.Println("Id2 error!")
	}

}
