package utils

import (
	"net"
)

type IKcpRequest interface {
	ICachedData;
	GetLEN()uint16;
	GetUID()uint32;
	GetRID()uint32;
	GetALL()[]byte;
	GetMsgBody()[]byte;
}
type IUdpRequest interface{
	IKcpRequest;
	GetAdr()net.Addr;
}
type IKcpResponse interface {
	ICachedData;
	IsBroadcast()bool;
	GetUID()uint32;
	GetSendData()[]byte;
	GetAllBDY()[]byte;
}
type IUdpResponse interface{
	ICachedData;
	IsBroadcast()bool;
	GetUID()uint32;
	GetSendData()[]byte;
	GetAllBDY()[]byte;
}