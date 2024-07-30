package Handles

import (
	"Server/Net"
	"Server/NetInterface"
	"fmt"
)

type KeepAliveHandler struct {
	Net.Handler
}

func (this *KeepAliveHandler) Handle(request NetInterface.IRequest) {
	Id2 := request.GetMsgID2()
	uid := request.GetMsgID3()
	switch Id2 {
	case 0:
		{
			gest := request.GetGesture()
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 0, uid, []byte(""), gest)
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	case 1:
		{
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 1, uid, []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	case 2:
		{
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 2, uid, []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	case 3:
		{
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 3, uid, []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	case 4:
		{
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 4, uid, []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	case 5:
		{
			if request.GetConnection().GetTeaConn() != nil {
				err := request.GetConnection().GetTeaConn().SendMsg(1, 5, uid, []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
				if err != nil {
					fmt.Println("Send msd error: ", err)
				}
			}
			break
		}
	default:
		{
			fmt.Println("Shouldn't be here in keepaliveHandle!")
			break
		}
	}

}
