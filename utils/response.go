package utils

import (
	"net"
)

type KcpRes struct{
	ICachedData;
	Broadcast bool;
	UID uint32;
	LEN uint16;
	BDY []byte;
}
func (me *KcpRes)IsBroadcast()bool{
	return me.Broadcast;
}
func (me *KcpRes)GetUID()uint32{
	return me.UID;
}
func (me *KcpRes)GetSendData()[]byte{
	return me.BDY[:me.LEN];
}

type UdpRes struct{
	ICachedData;
	Broadcast bool;
	ADR net.Addr;
	UID uint32;
	LEN uint16;
	BDY []byte;
}
func (me *UdpRes)IsBroadcast()bool{
	return me.Broadcast;
}
func (me *UdpRes)GetUID()uint32{
	return me.UID;
}
func (me *UdpRes)GetSendData()[]byte{
	return me.BDY[:me.LEN];
}
