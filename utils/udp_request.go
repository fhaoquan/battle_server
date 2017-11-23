package utils

import (
	"net"
	"encoding/binary"
)

/*
________________________________________________
udp_header|check_flag|usr_id|room_id|len|message|
____4_____|_____2____|___4__|___4___|_2_|_______|

*/

type udp_req struct {
	ICachedData;
	ADR net.Addr;
	Data []byte;
}
func (me *udp_req)Protocol()uint8{
	return Protocol_UDP;
}
func (me *udp_req)GetAdr()net.Addr{
	return me.ADR;
}
func (me *udp_req)OnReturn(){
	binary.LittleEndian.PutUint16(me.Data,0);
}
func (me *udp_req)Check()bool{
	return binary.LittleEndian.Uint16(me.Data[4:6])==12345;
}
func (me *udp_req)GetALL()[]byte{
	return me.Data;
}
func (me *udp_req)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.Data[14:16]);
}
func (me *udp_req)GetUID()uint32{
	return binary.LittleEndian.Uint32(me.Data[6:10]);
}
func (me *udp_req)GetRID()uint32{
	return binary.LittleEndian.Uint32(me.Data[10:14]);
}
func (me *udp_req)GetMsgBody()[]byte{
	return me.Data[16:me.GetLEN()+16];
}
func NewUdpReqPool(size int)(*MemoryPool){
	return NewMemoryPool(size, func(impl ICachedData) ICachedData {
		r:=&udp_req{
			impl,
			nil,
			make([]byte,MaxPktSize),
		}
		return r;
	});
}
