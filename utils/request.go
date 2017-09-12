package utils

import (
	"net"
	"encoding/binary"
)

type KcpReq struct {
	ICachedData;
	Data []byte;
}
func (me *KcpReq)OnReturn(){
	binary.LittleEndian.PutUint16(me.Data,0);
}
func (me *KcpReq)Check()bool{
	return binary.LittleEndian.Uint16(me.Data[0:2])==12345;
}
func (me *KcpReq)GetALL()[]byte{
	return me.Data;
}
func (me *KcpReq)GetLEN()uint16{
	return binary.LittleEndian.Uint16(me.Data[2:]);
}
func (me *KcpReq)GetUID()uint32{
	return binary.LittleEndian.Uint32(me.Data[4:]);
}
func (me *KcpReq)GetRID()uint32{
	return binary.LittleEndian.Uint32(me.Data[8:]);
}
func (me *KcpReq)GetMsgBody()[]byte{
	return me.Data[8:me.GetLEN()-8];
}

type UdpReq struct {
	ADR net.Addr;
	*KcpReq;
}
func (me *UdpReq)GetAdr()net.Addr{
	return me.ADR;
}

