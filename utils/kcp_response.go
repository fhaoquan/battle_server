package utils

import (
	"encoding/binary"
)

/*
_____________
|len|message|
|_2_|_______|

*/
type kcp_res struct{
	ICachedData;
	Broadcast bool;
	UID uint32;
	BDY []byte;
}
func (me *kcp_res)Protocol()uint8{
	return Protocol_KCP;
}
func (me *kcp_res)SetBroadcast(v bool){
	me.Broadcast=v;
}
func (me *kcp_res)SetUID(v uint32){
	me.UID=v;
}
func (me *kcp_res)GetBroadcast()bool{
	return me.Broadcast;
}
func (me *kcp_res)GetUID()uint32{
	return me.UID;
}
func (me *kcp_res)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.BDY);
}
func (me *kcp_res)GetSendData()[]byte{
	return me.BDY[:me.GetLEN()+2];
}
func (me *kcp_res)GetWriteBuffer()[]byte{
	return me.BDY;
}
func NewKcpResPool(size int)(*MemoryPool){
	return NewMemoryPool(size, func(impl ICachedData) ICachedData {
		return &kcp_res{
			impl,
			false,
			0,
			make([]byte,MaxPktSize),
		}
	});
}


