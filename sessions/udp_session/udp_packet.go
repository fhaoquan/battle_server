package udp_session

import "net"

type udp_packet struct {
	a net.Addr;
	l uint16;
	u uint32;
	r uint32;
	b []byte;
}
func (me *udp_packet)Clear(){
	me.a=nil;
	me.l=0;
	me.u=0;
	me.r=0;
}
func (me *udp_packet)GetLEN()uint16{
	return me.l;
}
func (me *udp_packet)GetUID()uint32{
	return me.u;
}
func (me *udp_packet)GetRID()uint32{
	return me.r;
}
func (me *udp_packet)GetRecvData()[]byte{
	return me.b[10:me.l];
}
func (me *udp_packet)GetADR()net.Addr{
	return me.a;
}