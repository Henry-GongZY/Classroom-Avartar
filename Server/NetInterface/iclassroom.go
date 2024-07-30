package NetInterface

type IClassroom interface {
	GetTeacher() IConnection

	AddStudent(IConnection, uint32) error

	DelStudent(uint32) error

	GetStudent(uint32) IConnection

	GetStudents() map[uint32]IConnection

	SetLessonid(uint32)

	Getlessonid() uint32

	Signin()

	GenStulist([]uint32)

	SetLessonname(string)

	Getlessonname() string
}
