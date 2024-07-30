package NetInterface

import "Server/ManagementInterface"

type IServer interface {
	Start()

	Stop()

	Serve()

	AddHandler(uint8, IHandler)

	GetConnMg() IConnManager

	GetDbCursor() ManagementInterface.IDatabase

	GetFilehandler() ManagementInterface.IFiletask
}
