package Net

import (
	"log"
	"strconv"
)

type Message struct {
	//消息ID
	id  uint8
	id2 uint8
	id3 string
	//数据
	data []byte
	//长度
	length int64
	//状态
	roll        float32
	pitch       float32
	yaw         float32
	min_ear     float32
	mar         float32
	mdst        float32
	LFronterArm float32
	LUpperArm   float32
	RFronterArm float32
	RUpperArm   float32
	//简易状态
	present bool
	handsup bool
	//额外数据
	mesg string
}

func (m *Message) SetGesture(gest []float32) {
	m.roll = gest[0]
	m.pitch = gest[1]
	m.yaw = gest[2]
	m.min_ear = gest[3]
	m.mar = gest[4]
	m.mdst = gest[5]
	m.LFronterArm = gest[6]
	m.LUpperArm = gest[7]
	m.RFronterArm = gest[8]
	m.RUpperArm = gest[9]
}

func (m *Message) GetGesture() []float32 {
	gest := make([]float32, 10)
	gest[0] = m.roll
	gest[1] = m.pitch
	gest[2] = m.yaw
	gest[3] = m.min_ear
	gest[4] = m.mar
	gest[5] = m.mdst
	gest[6] = m.LFronterArm
	gest[7] = m.LUpperArm
	gest[8] = m.RFronterArm
	gest[9] = m.RUpperArm
	log.Println(strconv.Itoa(int(gest[6])) + " " + strconv.Itoa(int(gest[7])) + " " + strconv.Itoa(int(gest[8])) + " " + strconv.Itoa(int(gest[9])) + " ")
	return gest
}

func (m *Message) GetMsgId() uint8 {
	return m.id
}

func (m *Message) GetMsgId2() uint8 {
	return m.id2
}

func (m *Message) GetMsgId3() string {
	return m.id3
}

func (m *Message) GetData() []byte {
	return m.data
}

func (m *Message) GetMsgLen() int64 {
	return m.length
}

func (m *Message) SetMesg(mesg string) {
	m.mesg = mesg
}

func (m *Message) GetMesg() string {
	return m.mesg
}

func (m *Message) SetMsgId(id uint8) {
	m.id = id
}

func (m *Message) SetMsgId2(id uint8) {
	m.id2 = id
}

func (m *Message) SetMsgId3(id string) {
	m.id3 = id
}

func (m *Message) SetData(data []byte) {
	m.data = data
}

func (m *Message) SetMsgLen(length int64) {
	m.length = length
}

//创建message包
func NewMsgPack(id uint8, id2 uint8, id3 string, data []byte, gest []float32) *Message {
	return &Message{
		id:          id,
		id2:         id2,
		id3:         id3,
		data:        data,
		roll:        gest[0],
		pitch:       gest[1],
		yaw:         gest[2],
		min_ear:     gest[3],
		mar:         gest[4],
		mdst:        gest[5],
		LFronterArm: gest[6],
		LUpperArm:   gest[7],
		RFronterArm: gest[8],
		RUpperArm:   gest[9],
		length:      0,
	}
}
