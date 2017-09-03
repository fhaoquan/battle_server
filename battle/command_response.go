package battle

import (
	"../utils"
	"net"
)

type kcp_response struct{
	utils.ICachedData;
	broadcast bool;
	uid uint32;
	len uint16;
	bdy []byte;
}
func (me *kcp_response)IsBroadcast()bool{
	return me.broadcast;
}
func (me *kcp_response)GetUID()uint32{
	return me.uid;
}
func (me *kcp_response)GetSendData()[]byte{
	return me.bdy[:me.len];
}

type udp_response struct{
	kcp_response;
	adr net.Addr;
}