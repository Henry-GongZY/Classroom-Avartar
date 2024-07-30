package Net

import (
	"Server/NetInterface"
	"Server/pb"
	"Server/utils"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
)

//TLV格式封包
type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取包头部长度
func (d *DataPack) GetHeadlen() uint32 {
	//"Length:0000",总共11字节
	return 11
}

//封包方法
func (d *DataPack) Pack(msg NetInterface.IMessage) ([]byte, error) {
	//创建缓冲dataBuff
	dataBuff := &pb.Packet{}

	//打包
	//写入MsgId
	if msg.GetMsgId() <= 10 {
		dataBuff.Id1 = pb.Id1(msg.GetMsgId())
	} else {
		return nil, errors.New("Protobuf Id1 error!")
	}

	//写入MsgId2
	dataBuff.Id2 = uint32(msg.GetMsgId2())

	//写入MsgId3
	dataBuff.SorTid = msg.GetMsgId3()

	//写入Data
	switch dataBuff.Id1 {
	case pb.Id1_KEEPALIVE:
		{
			//坐标更新
			gest := msg.GetGesture()
			dataBuff.Gesture = &pb.Gesture{}
			dataBuff.Gesture.Roll = gest[0]
			dataBuff.Gesture.Pitch = gest[1]
			dataBuff.Gesture.Yaw = gest[2]
			dataBuff.Gesture.MinEar = gest[3]
			dataBuff.Gesture.Mar = gest[4]
			dataBuff.Gesture.Mdst = gest[5]
			dataBuff.Gesture.LFronterArm = gest[6]
			dataBuff.Gesture.LUpperArm = gest[7]
			dataBuff.Gesture.RFronterArm = gest[8]
			dataBuff.Gesture.RUpperArm = gest[9]
			//log.Println(strconv.Itoa(int(gest[6])) + " " + strconv.Itoa(int(gest[7])) + " " + strconv.Itoa(int(gest[8])) + " " + strconv.Itoa(int(gest[9])) + " ")
			break
		}
	case pb.Id1_FILE:
		{
			if dataBuff.Id2 == 3 {
				dataBuff.Mesg = string(msg.GetData())
			} else {
				dataBuff.Filedata = msg.GetData()
			}
			break
		}
	case pb.Id1_LESSONS:
		{
			dataBuff.Filedata = msg.GetData()
			break
		}
	case pb.Id1_FILEUPLOAD:
		{
			dataBuff.Filedata = msg.GetData()
		}
	default:
		dataBuff.Mesg = string(msg.GetData())
	}

	return proto.Marshal(dataBuff)
}

//拆包方法
func (d *DataPack) Unpack(binaryData []byte) (NetInterface.IMessage, error) {
	dataBuff := &pb.Packet{}
	err := proto.Unmarshal(binaryData, dataBuff)
	if err != nil {
		return nil, err
	}

	//保存信息,创建空包
	msg := NewMsgPack(0, 0, "0", []byte(""), []float32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0})

	//读取MsgId
	MsgId1 := uint8(dataBuff.GetId1())
	msg.SetMsgId(MsgId1)

	//读取MsgId2
	MsgId2 := uint8(dataBuff.GetId2())
	msg.SetMsgId2(MsgId2)

	//读取MsgId3
	MsgId3 := dataBuff.GetSorTid()
	msg.SetMsgId3(MsgId3)

	//读取Data
	//读取登陆密码
	switch MsgId1 {
	case 0:
		{
			msg.SetData([]byte(dataBuff.GetMesg()))
			break
		}
	case 1:
		{
			if MsgId2 == 0 {
				msg.SetGesture([]float32{dataBuff.Gesture.Roll,
					dataBuff.Gesture.Pitch,
					dataBuff.Gesture.Yaw,
					dataBuff.Gesture.MinEar,
					dataBuff.Gesture.Mar,
					dataBuff.Gesture.Mdst,
					dataBuff.Gesture.LFronterArm,
					dataBuff.Gesture.LUpperArm,
					dataBuff.Gesture.RFronterArm,
					dataBuff.Gesture.RUpperArm,
				})
				log.Println("解包" + strconv.Itoa(int(dataBuff.Gesture.LFronterArm)) + " " + strconv.Itoa(int(dataBuff.Gesture.LUpperArm)) + " " + strconv.Itoa(int(dataBuff.Gesture.RFronterArm)) + " " + strconv.Itoa(int(dataBuff.Gesture.RUpperArm)) + " ")
			}
			break
		}
	case 3:
		{
			msg.SetMesg(dataBuff.GetMesg())
			break
		}
	case 4:
		{
			msg.SetMesg(dataBuff.GetMesg())
			if MsgId2 == 2 {
				msg.SetData(dataBuff.GetFiledata())
			}
			break
		}
	default:
		break
	}

	//判断包长，丢弃超长的包
	if utils.GlobalObject.MaxPackageSize > 0 && msg.GetMsgLen() > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New(fmt.Sprintf("Too Large msg data recv: ", msg.GetMsgLen()))
	}

	return msg, nil
}
