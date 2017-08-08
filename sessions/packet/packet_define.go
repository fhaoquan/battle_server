package packet

import (
	"../../utils"
	"net"
)

type IKcpRequest interface {
	utils.ICachedData;
	ReadAt(net.Conn)error;
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetBDY()[]byte;
}
type IUdpRequest interface{
	utils.ICachedData;
	ReadAt(net.PacketConn)error;
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetBDY()[]byte;
	GetAdr()net.Addr;
}
type IKcpResponse interface {
	utils.ICachedData;
	IsBroadcast()bool;
	GetUID()uint32;
	GetBDY()[]byte;
}
type IUdpResponse interface{
	IKcpResponse;
	GetAdr()net.Addr;
}