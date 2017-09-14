package utils

import (
	"encoding/binary"
)
const CMD_pingpong =uint8(0);
const CMD_battle_wating_start =uint8(1);
const CMD_battle_all =uint8(2);
const CMD_unit_movment =uint8(3);
const CMD_create_unit =uint8(4);
const CMD_attack_done =uint8(5);
const CMD_attack_start =uint8(6);
const CMD_battle_start =uint8(7);

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
