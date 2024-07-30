package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Login struct {
	Stuid  uint32 `json:"student_id"`
	Signed bool   `json:"signin"`
}

type SigninObj struct {
	TimeStamp string  `json:"time_stamp"`
	Students  []Login `json:"student_login"`
}

type SigninObjs struct {
	Singlelesson []SigninObj `json:"student_situation"`
}

func (s *SigninObj) Gensignin(signin map[uint32]bool, lessonname string) {
	//读取签到json表
	data, _ := ioutil.ReadFile("Files/" + lessonname + "/login.json")
	fd, _ := os.OpenFile("Files/"+lessonname+"/login.json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0766)
	defer fd.Close()

	//记录登录情况
	var sg SigninObjs
	if len(data) != 0 {
		err := json.Unmarshal(data, &sg)
		if err != nil {
			fmt.Println("Unmarshal error: ", err)
		}
	}

	for k, v := range signin {
		s.Students = append(s.Students, Login{k, v})
	}
	t := time.Now().Format("2006-01-02 15:04:05")
	s.TimeStamp = strings.Split(t, " ")[0]
	sg.Singlelesson = append(sg.Singlelesson, *s)

	//格式化
	b, err := json.Marshal(sg)
	if err != nil {
		fmt.Println("Marshall json error:", err)
	}

	//存储为json文件
	_, err = fd.Write(b)
	if err != nil {
		fmt.Println("Write json error: ", err)
	}
}
