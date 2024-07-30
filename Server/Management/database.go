package Management

import (
	"Server/Data"
	"Server/ManagementInterface"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//数据库查询操作
type Database struct {
	db *sqlx.DB
}

func NewDatabase() ManagementInterface.IDatabase {
	d := &Database{}
	return d
}

//连接数据库
func (d *Database) Connection() {
	database, err := sqlx.Open("mysql", "root@tcp(127.0.0.1:3306)/gra_design")
	if err != nil {
		fmt.Println("Mysql service connection error: ", err)
		return
	}
	d.db = database
}

//查找学生
func (d *Database) FindStu(stuid uint32) []Data.Student {
	var stu []Data.Student
	querySQL := "select Stuid,Stuname from student where Stuid = ?"
	if err := d.db.Select(&stu, querySQL, stuid); err != nil {
		fmt.Println("Mysql find error: ", err)
	}
	return stu
}

//查找老师
func (d *Database) FindTea(stuid uint32) []Data.Teacher {
	var teacher []Data.Teacher
	querySQL := "select Tid,Tname from teacher where Tid = ?"
	if err := d.db.Select(&teacher, querySQL, stuid); err != nil {
		fmt.Println("Mysql find error: ", err)
	}
	return teacher
}

//学生登录
func (d *Database) StuLogin(stuid uint32, password string) (bool, string) {
	str := ""
	querySQL := "select Stuname from student where Stuid = ? and Password = ?"
	if err := d.db.Get(&str, querySQL, stuid, password); err != nil {
		fmt.Println("Mysql login error: ", err)
	}

	if len(str) == 0 {
		return false, str
	}
	return true, str
}

//老师登录
func (d *Database) TeaLogin(teacherid uint32, password string) (bool, string) {
	str := ""
	querySQL := "select Tname from teacher where Tid = ? and Password = ?"
	if err := d.db.Get(&str, querySQL, teacherid, password); err != nil {
		fmt.Println("Mysql login error: ", err)
	}
	if len(str) == 0 {
		return false, str
	}
	return true, str
}

//判断老师现在是否有课程
func (d *Database) TeaClass(teacherid uint32) (uint32, string, bool) {
	var str []Data.Lesson

	//测试用
	querySQL := "select Lid,Lname from teaching natural join lesson where Tid = ?"
	if err := d.db.Select(&str, querySQL, teacherid); err != nil {
		fmt.Println("Mysql error: ", err)
	}
	//正常使用
	//querySQL := "select Lid,Lname from teaching natural join lesson where Tid = ? and DAY = ? and STARTTIME <= ? and ENDTIME > ?"
	//weekday := int(time.Now().Weekday())
	//nowtime := strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second())
	//if err := d.db.Select(&str, querySQL, teacherid, weekday, nowtime, nowtime); err != nil {
	//	fmt.Println("Mysql error: ", err)
	//}

	if len(str) != 0 {
		return str[0].Lessonid, str[0].Lessonname, true
	} else {
		return 0, "", false
	}
}

//判断学生现在是否有课程
func (d *Database) StuClass(studentid uint32) (uint32, string, bool) {
	var str []Data.Lesson

	//测试用
	querySQL := "select Lid,Lname from lesson natural join participant where Stuid = ?"
	if err := d.db.Select(&str, querySQL, studentid); err != nil {
		fmt.Println("Mysql error: ", err)
	}
	//正常使用
	//querySQL := "select Lid,Lname from lesson natural join participant where Stuid = ? and DAY = ? and STARTTIME <= ? and ENDTIME > ?"
	//weekday := int(time.Now().Weekday())
	//nowtime := strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute()) + ":" + strconv.Itoa(time.Now().Second())
	//if err := d.db.Select(&str, querySQL, studentid, weekday, nowtime, nowtime); err != nil {
	//	fmt.Println("Mysql error: ", err)
	//}

	if len(str) != 0 {
		return str[0].Lessonid, str[0].Lessonname, true
	} else {
		return 0, "", false
	}
}

//查找特定学生参与的课程
func (d *Database) StuClasses(studentid uint32) ([]Data.Lesson, bool) {
	var str []Data.Lesson

	querySQL := "select Lid,Lname from lesson natural join participant where Stuid = ?"
	if err := d.db.Select(&str, querySQL, studentid); err != nil {
		fmt.Println("Mysql error: ", err)
	}

	if len(str) != 0 {
		return str, true
	} else {
		return nil, false
	}
}

//查找特定老师带领的课程
func (d *Database) TeaClasses(teacherid uint32) ([]Data.Lesson, bool) {
	var str []Data.Lesson

	querySQL := "select Lid,Lname from lesson natural join teaching where Tid = ?"
	if err := d.db.Select(&str, querySQL, teacherid); err != nil {
		fmt.Println("Mysql error: ", err)
	}

	if len(str) != 0 {
		return str, true
	} else {
		return nil, false
	}
}

//查找参与课程的所有学生
func (d *Database) Students(Lessonid uint32) []uint32 {
	var students []uint32

	querySQL := "select Stuid from participant where Lid = ?"
	if err := d.db.Select(&students, querySQL, Lessonid); err != nil {
		fmt.Println("Mysql error: ", err)
	}

	return students
}
