package utils

import (
	"net"
)

type IKcpRequest interface {
	ICachedData;
	ReadAt(net.Conn)error;
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetRecvData()[]byte;
}
type IUdpRequest interface{
	ICachedData;
	ReadAt(net.PacketConn)error;
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetRecvData()[]byte;
	GetAdr()net.Addr;
}
type IKcpResponse interface {
	ICachedData;
	IsBroadcast()bool;
	GetUID()uint32;
	SetUID(uint32);
	GetSendData()[]byte;
	GetAllBDY()[]byte;
}
type IUdpResponse interface{
	ICachedData;
	IsBroadcast()bool;
	GetUID()uint32;
	SetUID(uint32);
	GetSendData()[]byte;
	GetAllBDY()[]byte;
}