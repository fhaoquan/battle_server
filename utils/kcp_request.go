package utils

import (
	"encoding/binary"
)

/*
______________________________________
check_flag|usr_id|room_id|len|message|
_____2____|___4__|___4___|_2_|_______|

*/
type kcp_req struct {
	ICachedData;
	Data []byte;
}
func (me *kcp_req)Protocol()uint8{
	return Protocol_KCP;
}
func (me *kcp_req)OnReturn(){
	binary.LittleEndian.PutUint16(me.Data,0);
}
func (me *kcp_req)Check()bool{
	return binary.LittleEndian.Uint16(me.Data[0:2])==12345;
}
func (me *kcp_req)GetALL()[]byte{
	return me.Data;
}
func (me *kcp_req)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.Data[10:12]);
}
func (me *kcp_req)GetUID()uint32{
	return binary.LittleEndian.Uint32(me.Data[2:6]);
}
func (me *kcp_req)GetRID()uint32{
	return binary.LittleEndian.Uint32(me.Data[6:10]);
}
func (me *kcp_req)GetMsgBody()[]byte{
	return me.Data[12:me.GetLEN()+12];
}

func NewKcpReqPool(size int)(*MemoryPool){
	return NewMemoryPool(size, func(impl ICachedData) ICachedData {
		r:=&kcp_req{
			impl,
			make([]byte,MaxPktSize),
		}
		return r;
	});
}

