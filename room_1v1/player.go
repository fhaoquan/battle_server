package room_1v1

import (
	"../gateway"
)

type player struct {
	uid				uint32
	name			string
	session			gateway.Session;
}
func (me *player)send_kcp(data []byte){
	if me.session==nil{
		return ;
	}
	defer func() {
		if e:=recover();e!=nil{
			me.session.Close();
		}
	}()
	me.session.KcpSend(data);
}
func (me *player)send_udp(data []byte){
	if me.session==nil{
		return ;
	}
	defer func() {
		if e:=recover();e!=nil{
			me.session.Close();
		}
	}()
	me.session.UdpSend(data);
}
