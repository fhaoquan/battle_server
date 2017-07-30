package kcp_session

type kcp_packet struct {
	l uint16;
	u uint32;
	r uint32;
	b []byte;
}
func (me *kcp_packet)Clear(){
	me.l=0;
	me.u=0;
	me.r=0;
}
func (me *kcp_packet)GetLEN()uint16{
	return me.l;
}
func (me *kcp_packet)GetUID()uint32{
	return me.u;
}
func (me *kcp_packet)GetRID()uint32{
	return me.r;
}
func (me *kcp_packet)GetBDY()[]byte{
	return me.b[10:me.l];
}