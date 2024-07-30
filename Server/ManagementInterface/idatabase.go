package ManagementInterface

import (
	"Server/Data"
)

type IDatabase interface {
	Connection()

	FindStu(uint32) []Data.Student

	FindTea(uint32) []Data.Teacher

	StuLogin(uint32, string) (bool, string)

	TeaLogin(uint32, string) (bool, string)

	TeaClass(uint32) (uint32, string, bool)

	StuClass(uint32) (uint32, string, bool)

	TeaClasses(uint32) ([]Data.Lesson, bool)

	StuClasses(uint32) ([]Data.Lesson, bool)

	Students(uint32) []uint32
}
