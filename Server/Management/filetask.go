package Management

import (
	"Server/Data"
	"fmt"
	"os"
	"strings"
)

type Filetask struct {
	//文件
	Filemap map[string]*os.File
	//文件管道
	Filechannel chan Data.FileCommand
}

func (f *Filetask) Run() {
	f.Filemap = make(map[string]*os.File)
	go f.File()
}

func (f *Filetask) File() {
	for {
		select {
		case a := <-f.Filechannel:
			if a.Command == 1 {
				f.New(a.Name)
			} else if a.Command == 2 {
				f.Write(a.Name, a.Data)
			} else if a.Command == 3 {
				f.Close(a.Name)
			}
		}
	}
}

func (f *Filetask) New(filename string) {
	fmt.Println(filename)
	//文件名+课程名+学生id
	ff := strings.Split(filename, "/")

	//新建课程文件夹
	if _, err := os.Stat("./Files/" + ff[1]); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("./Files/"+ff[1], 0777)
			if err != nil {
				fmt.Println("Mkdir error: ", err)
			}
		}
	}

	//新建学生文件夹
	if _, err := os.Stat("./Files/" + ff[1] + "/" + ff[2]); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("./Files/"+ff[1]+"/"+ff[2], 0777)
			if err != nil {
				fmt.Println("Mkdir error: ", err)
			}
		}
	}

	file, err := os.Create("./Files/" + ff[1] + "/" + ff[2] + "/" + ff[0])
	if err != nil {
		fmt.Println("Create File error: ", err)
	}
	f.Filemap[filename] = file
}

func (f *Filetask) Write(filename string, data []byte) {
	_, err := f.Filemap[filename].Write(data)
	if err != nil {
		fmt.Println("Write File error: ", err)
	}
}

func (f *Filetask) Close(filename string) {
	f.Filemap[filename].Close()
	delete(f.Filemap, filename)
}

func (f *Filetask) GetChannel() chan Data.FileCommand {
	return f.Filechannel
}
