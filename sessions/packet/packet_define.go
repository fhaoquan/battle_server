package packet

import "net"

type IKcpPacket interface {
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetBDY()[]byte;
}
type IUdpPacket interface{
	IKcpPacket;
	GetAdr()net.Addr;
}

type PacketChannel struct{

}
