package Data

type Teacher struct {
	Teacherid   uint32 `db:"Tid"`
	Teachername string `db:"Tname"`
	Password    string `db:"Password"`
}

type Student struct {
	Stuid    uint32 `db:"Stuid"`
	Stuname  string `db:"Stuname"`
	Password string `db:"Password"`
}

type Lesson struct {
	Lessonid   uint32 `db:"Lid"`
	Lessonname string `db:"Lname"`
}

type FileCommand struct {
	Data    []byte
	Name    string
	Command int
}
