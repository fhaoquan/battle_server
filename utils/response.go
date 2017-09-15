package utils

import (
	"encoding/binary"
)

type KcpRes struct{
	ICachedData;
	Broadcast bool;
	UID uint32;
	BDY []byte;
}
func (me *KcpRes)IsBroadcast()bool{
	return me.Broadcast;
}
func (me *KcpRes)GetUID()uint32{
	return me.UID;
}
func (me *KcpRes)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.BDY);
}
func (me *KcpRes)GetSendData()[]byte{
	return me.BDY[:me.GetLEN()+2];
}
func (me *KcpRes)GetAllBDY()[]byte{
	return me.BDY;
}

type UdpRes struct{
	*KcpRes;
};
