package room

import (
	"../utils"
	"net"
)

type recv_channel struct{
	kcp_chan chan utils.IDataOwner;
	kcp_pool *utils.MemoryPool;
	udp_chan chan utils.IDataOwner;
	udp_pool *utils.MemoryPool;
}

func (r *recv_channel)OnKcp(len uint16,uid uint32,rid uint32,bdy []byte)(e error){
	o:=r.kcp_pool.PopOne();
	o.GetUserData().(*kcp_packet).set_all(len,uid,rid,bdy);
	r.kcp_chan<-o;
	return nil;
}
func (r *recv_channel)OnUdp(adr net.Addr,len uint16,uid uint32,rid uint32,bdy []byte)(e error){
	o:=r.udp_pool.PopOne();
	o.GetUserData().(*udp_packet).set_all(adr,len,uid,rid,bdy);
	r.udp_chan<-o;
	return nil;
}

func new_recv_channel()(*recv_channel){
	size:=16;
	return &recv_channel{
		kcp_chan:make(chan utils.IDataOwner,size),
		kcp_pool:utils.NewMemoryPool(size, func()interface{} {
			return &kcp_packet{0,0,0,make([]byte,utils.MaxPktSize)};
		}),
		udp_chan:make(chan utils.IDataOwner,size),
		udp_pool:utils.NewMemoryPool(size, func()interface{} {
			return &udp_packet{nil,&kcp_packet{0,0,0,make([]byte,utils.MaxPktSize)}};
		}),
	}
}