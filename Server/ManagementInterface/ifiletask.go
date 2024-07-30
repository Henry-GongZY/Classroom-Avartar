package ManagementInterface

import (
	"Server/Data"
)

type IFiletask interface {
	Run()

	File()

	New(string)

	Write(string, []byte)

	Close(string)

	GetChannel() chan Data.FileCommand
}
