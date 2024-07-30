package Handles

import (
	"Server/Net"
	"Server/NetInterface"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type FileHandler struct {
	Net.Handler
	Channel chan bool
}

func (this *FileHandler) PreHandle(request NetInterface.IRequest) {
	this.Channel = make(chan bool, 1)
}

func (this *FileHandler) Handle(request NetInterface.IRequest) {
	Id2 := request.GetMsgID2()
	SorTid := request.GetMsgID3()
	mesg := request.GetMesg()
	switch Id2 {
	case 0:
		{
			strdata := strings.Split(mesg, "/")

			filename := strings.Split(strdata[0], " ")[0]
			lesson := strdata[1]

			f, err := os.Open("./Files/" + lesson + "/" + SorTid + "/" + filename)
			if err != nil {
				fmt.Println("os.Open err = ", err)
				return
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("File close error!")
				}
			}(f)

			if err := request.GetConnection().SendMsg(3, 3, "0", []byte(filename), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
				fmt.Println("Filename send error!")
			}

			time.Sleep(1 * time.Second)

			buf := make([]byte, 1024*4)
			//读文件内容，读多少发多少
			for {
				n, err := f.Read(buf) //从文件读取内容
				if err != nil {
					if err == io.EOF {
						fmt.Println("File transportation finished!")
					} else {
						fmt.Println("f.Read err = ", err)
					}
					this.Channel <- true
					return
				}
				//发送内容
				if err := request.GetConnection().SendMsg(3, 2, "0", buf[:n], []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
					fmt.Println("File write error!")
				}
			}
		}
	case 1:
		{
			strdata := strings.Split(mesg, "/")

			filename := strings.Split(strdata[0], " ")[0]
			lesson := strdata[1]

			f, err := os.Open("./Files/" + lesson + "/Teacher/" + filename)
			if err != nil {
				fmt.Println("os.Open err = ", err)
				return
			}

			defer func(f *os.File) {
				err := f.Close()
				if err != nil {
					fmt.Println("File close error!")
				}
			}(f)

			if err := request.GetConnection().SendMsg(3, 3, "0", []byte(filename), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
				fmt.Println("Filename send error!")
			}

			time.Sleep(1 * time.Second)

			buf := make([]byte, 1024*4)
			//读文件内容，读多少发多少
			for {
				n, err := f.Read(buf) //从文件读取内容
				if err != nil {
					if err == io.EOF {
						fmt.Println("File transportation finished!")
					} else {
						fmt.Println("f.Read err = ", err)
					}
					this.Channel <- true
					return
				}
				//发送内容
				if err := request.GetConnection().SendMsg(3, 2, "0", buf[:n], []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
					fmt.Println("File write error!")
				}
			}
		}
	}

}

func (this *FileHandler) PostHandle(request NetInterface.IRequest) {
	//阻塞接受
	for {
		select {
		//文件发送完成就分开
		case <-this.Channel:
			if err := request.GetConnection().SendMsg(3, 1, "0", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}); err != nil {
				fmt.Println("Goodbye error!")
			}
			close(this.Channel)
			return
		}
	}
}
