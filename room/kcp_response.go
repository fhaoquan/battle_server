package room

import "../utils"


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
func (me *kcp_response)SetBroadcast(b bool){
	me.broadcast=b;
}
func (me *kcp_response)SetUID(u uint32){
	me.uid=u;
}
func (me *kcp_response)GetAllBDY()[]byte{
	return me.bdy;
}
