package main

import (
	"Server/Handles"
	"Server/Management"
	"Server/Net"
)

func main() {
	db := Management.NewDatabase()
	db.Connection()

	s := Net.NewServer(db)
	s.AddHandler(0, &Handles.LoginHandler{})
	s.AddHandler(1, &Handles.KeepAliveHandler{})
	s.AddHandler(2, &Handles.InfoHandler{})
	s.AddHandler(3, &Handles.FileHandler{Channel: make(chan bool, 1)})
	s.AddHandler(4, &Handles.FileUploadHandler{Handleon: false})
	s.AddHandler(5, &Handles.LsHandler{})
	s.AddHandler(9, &Handles.DebugHandler{})
	s.Serve()
}
