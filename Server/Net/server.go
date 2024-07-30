package Net

import (
	"Server/Data"
	"Server/Management"
	"Server/ManagementInterface"
	"Server/NetInterface"
	"Server/utils"
	"fmt"
	"net"
)

type Server struct {
	IPVersion string
	IPAddress string
	Port      int
	Parser    NetInterface.ITask
	ConnMg    NetInterface.IConnManager
	db        ManagementInterface.IDatabase
	fh        ManagementInterface.IFiletask
}

//启动函数
func (s *Server) Start() {
	go func() {
		//启动工作池
		s.Parser.StartWorkerPool()

		//启动文件管理
		s.fh = &Management.Filetask{Filechannel: make(chan Data.FileCommand, 10)}
		s.fh.Run()

		//分配IP地址
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IPAddress, s.Port))
		if err != nil {
			fmt.Println("Error at TCPAddr:", err)
		}

		Listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("Error at TCPListener:", err)
		}

		fmt.Println("Server started!")
		var cid uint32
		cid = 0

		for {
			fmt.Println("Waiting for connection!")
			conn, err := Listener.AcceptTCP()
			if err != nil {
				fmt.Println("Error at Acception:", err)
				continue
			}

			//判断链接是否超过最大个数，超过则关闭新的连接
			if s.ConnMg.Count() >= utils.GlobalObject.MaxConn {
				fmt.Println("Connection overflow!")
				err := conn.Close()
				if err != nil {
					return
				}
				continue
			}

			//链接绑定处理方法，启动连接业务
			dealConn := NewConnection(s, conn, cid, s.Parser)
			cid++

			go dealConn.Start()
		}
	}()
}

//终止函数
func (s *Server) Stop() {
	//释放和回收资源
	fmt.Println("Server stop!")
	s.ConnMg.Clear()
}

//服务事项
func (s *Server) Serve() {
	s.Start()
	//TODO:等待补充

	select {}
}

func (s *Server) AddHandler(msgId uint8, handler NetInterface.IHandler) {
	s.Parser.AddHandler(msgId, handler)
	fmt.Println("Add router succeed!")
}

func (s *Server) GetConnMg() NetInterface.IConnManager {
	return s.ConnMg
}

func (s *Server) GetDbCursor() ManagementInterface.IDatabase {
	return s.db
}

func (s *Server) GetFilehandler() ManagementInterface.IFiletask {
	return s.fh
}

//初始化
func NewServer(database ManagementInterface.IDatabase) NetInterface.IServer {
	s := &Server{
		IPVersion: utils.GlobalObject.GetVersion(),
		IPAddress: utils.GlobalObject.GetHost(),
		Port:      utils.GlobalObject.GetPort(),
		Parser:    NewParser(),
		ConnMg:    NewConnManager(),
		db:        database,
	}

	return s
}
