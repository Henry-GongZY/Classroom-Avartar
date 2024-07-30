package Handles

import (
	"Server/Net"
	"Server/NetInterface"
	"fmt"
)

type DebugHandler struct {
	Net.Handler
}

func (this *DebugHandler) Handle(request NetInterface.IRequest) {
	fmt.Println("Call Handler")
	fmt.Println("Recv from client: msgId = ", request.GetMsgID(),
		"msgId2 = ", request.GetMsgID2(), "msgId3 = ", request.GetMsgID3(), "Data:", string(request.GetData()))
	err := request.GetConnection().SendMsg(0, 0, "1811358", []byte("wcnm"), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return
	}
}
