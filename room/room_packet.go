package room

import "net"

type kcp_packet struct {
	len uint16;
	uid uint32;
	rid uint32;
	bdy []byte;
}
func (me *kcp_packet)set_all(len uint16,uid uint32,rid uint32,bdy []byte)(*kcp_packet){
	me.len=len;
	me.uid=uid;
	me.rid=rid;
	copy(me.bdy,bdy);
	return me;
}
type udp_packet struct {
	adr net.Addr;
	*kcp_packet;
}
func (me *udp_packet)set_all(adr net.Addr,len uint16,uid uint32,rid uint32,bdy []byte)(*udp_packet){
	me.adr=adr;
	me.kcp_packet.set_all(len,uid,rid,bdy)
	return me;
}
