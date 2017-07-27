package room

import (
	"net"
)

type recv_channel struct{
	kcp_chan chan func(func(uid uint32,rid uint32,bdy []byte)bool);
	udp_chan chan func(func(adr net.Addr,uid uint32,rid uint32,bdy []byte)bool);
}

func (r *recv_channel)OnKcp(f func(func(uid uint32,rid uint32,bdy []byte)bool))(e error){
	r.kcp_chan<-f;
	return nil;
}
func (r *recv_channel)OnUdp(f func(func(adr net.Addr,uid uint32,rid uint32,bdy []byte)bool))(e error){
	r.udp_chan<-f;
	return nil;
}

func new_recv_channel()(*recv_channel){
	return &recv_channel{
		make(chan func(func(uid uint32,rid uint32,bdy []byte)bool),16),
		make(chan func(func(adr net.Addr,uid uint32,rid uint32,bdy []byte)bool),16),
	}
}