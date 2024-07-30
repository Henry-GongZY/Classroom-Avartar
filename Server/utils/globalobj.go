package utils

import (
	"Server/NetInterface"
	"encoding/json"
	"io/ioutil"
)

type GlobalObj struct {
	//服务器
	TcpServer NetInterface.IServer
	//IP
	Host string
	//端口
	TcpPort int
	//版本号
	Version string
	//最大连接数目
	MaxConn int
	//包大小
	MaxPackageSize int64
	//工作池数量
	WorkerPoolSize uint32
	//消息队列长度
	MaxWorkerTaskLen uint32
}

var GlobalObject *GlobalObj

func init() {
	GlobalObject = &GlobalObj{
		Version:          "tcp4",
		TcpPort:          5665,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   512,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
	}

	GlobalObject.Reload()
}

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/serverconfig.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func (g *GlobalObj) GetVersion() string {
	return g.Version
}

func (g *GlobalObj) GetHost() string {
	return g.Host
}

func (g *GlobalObj) GetPort() int {
	return g.TcpPort
}
