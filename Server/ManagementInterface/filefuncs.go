package ManagementInterface

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func GetFileMd5(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("os Open error")
		return "", err
	}
	MD5 := md5.New()
	_, err = io.Copy(MD5, file)
	if err != nil {
		fmt.Println("io copy error")
		return "", err
	}
	md5Str := hex.EncodeToString(MD5.Sum(nil))
	return md5Str, nil
}

func GetTeaDoclist(lesson string, uid string, IsStu bool) []string {
	filemeta := []string{}
	if IsStu {
		//新建课程文件夹
		if _, err := os.Stat("./Files/" + lesson); err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir("./Files/"+lesson, 0777)
				if err != nil {
					fmt.Println("Mkdir error: ", err)
				}
			}
		}

		//新建教师文件夹
		if _, err := os.Stat("./Files/" + lesson + "/Teacher"); err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir("./Files/"+lesson+"/Teacher", 0777)
				if err != nil {
					fmt.Println("Mkdir error: ", err)
				}
			}
		}

		files, _ := ioutil.ReadDir("./Files/" + lesson + "/Teacher")

		for _, f := range files {
			year, month, day := f.ModTime().Date()
			hour, min, sec := f.ModTime().Clock()
			monthstr := month.String()
			Timestamp := fmt.Sprintf("%d-%s-%d %d:%d:%d",
				year, monthstr, day, hour, min, sec)
			filemeta = append(filemeta, f.Name()+" "+Timestamp)
		}
	} else {

	}

	return filemeta
}

func GetDoclist(lesson string, uid string, IsStu bool) []string {
	filemeta := []string{}
	if IsStu {
		//新建课程文件夹
		if _, err := os.Stat("./Files/" + lesson); err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir("./Files/"+lesson, 0777)
				if err != nil {
					fmt.Println("Mkdir error: ", err)
				}
			}
		}

		//新建学生文件夹
		if _, err := os.Stat("./Files/" + lesson + "/" + uid); err != nil {
			if os.IsNotExist(err) {
				err := os.Mkdir("./Files/"+lesson+"/"+uid, 0777)
				if err != nil {
					fmt.Println("Mkdir error: ", err)
				}
			}
		}

		files, _ := ioutil.ReadDir("./Files/" + lesson + "/" + uid)

		for _, f := range files {
			year, month, day := f.ModTime().Date()
			hour, min, sec := f.ModTime().Clock()
			monthstr := month.String()
			Timestamp := fmt.Sprintf("%d-%s-%d %d:%d:%d",
				year, monthstr, day, hour, min, sec)
			filemeta = append(filemeta, f.Name()+" "+Timestamp)
		}

	} else {

	}

	return filemeta
}

func DelDoclist(lesson string, uid string, IsStu bool, filename string) {
	if IsStu {
		err := os.Remove("./Files/" + lesson + "/" + uid + "/" + filename)
		if err != nil {
			fmt.Println(err)
		}
	} else {

	}
}
