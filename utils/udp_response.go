package utils

import "encoding/binary"

/*
______________________
udp_header|len|message|
____4_____|_2_|_______|

*/
type udp_res struct{
	ICachedData;
	Broadcast bool;
	UID uint32;
	BDY []byte;
};
func (me *udp_res)Protocol()uint8{
	return Protocol_UDP;
}
func (me *udp_res)SetBroadcast(v bool){
	me.Broadcast=v;
}
func (me *udp_res)SetUID(v uint32){
	me.UID=v
}
func (me *udp_res)GetBroadcast()bool{
	return me.Broadcast;
}
func (me *udp_res)GetUID()uint32{
	return me.UID;
}
func (me *udp_res)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.BDY[4:6]);
}
func (me *udp_res)GetSendData()[]byte{
	return me.BDY[:me.GetLEN()+6];
}
func (me *udp_res)GetWriteBuffer()[]byte{
	return me.BDY[4:];
}
func NewUdpResPool(size int)(*MemoryPool){
	return NewMemoryPool(size, func(impl ICachedData) ICachedData {
		r:=&udp_res{
			impl,
			false,
			0,
			make([]byte,MaxPktSize),
		}
		binary.LittleEndian.PutUint32(r.BDY,uint32(UdpPktHeader));
		return r;
	});
}
