package room

import (
	"time"
	"net"
	"../sessions/packet"
)

const MAX_CMD_ID  = 255;

type Room struct{
	id uint32;
	started bool;
	players map[uint32]*player;
	*recv_channel;
	cmd_handlers []func([]byte,*Room);
	timer_handlers []func(*Room);

}
func (r *Room)GetID()uint32{
	return r.id;
}
func (r *Room)SetID(v uint32)*Room{
	r.id=v;
	return r;
}
func (r *Room)Start()*Room{
	if(r.started){
		return r;
	}
	r.started=true;
	timer:=make(chan int,1);
	go func(){
		for{
			select {
			case dat:=<-r.kcp_chan:
				bdy:=dat.GetUserData().(*kcp_packet).bdy;
				if(r.cmd_handlers[bdy[0]]!=nil){
					r.cmd_handlers[bdy[0]](bdy[1:],r);
				}
				dat.Return();
			case dat:=<-r.udp_chan:
				bdy:=dat.GetUserData().(*udp_packet).bdy;
				if(r.cmd_handlers[bdy[0]]!=nil){
					r.cmd_handlers[bdy[0]](bdy[1:],r);
				}
				dat.Return();
			case tid:=<-timer:
				if(r.timer_handlers[tid]!=nil){
					r.timer_handlers[tid](r);
				}
			}
		}
	}();
	go func(){
		for{
			time.Sleep(time.Millisecond*30);
			timer<-1;
		}
	}();
	return r;
}
func NewRoom()(*S_room_builder){
	return &S_room_builder{
		&Room{
			0,
			false,
			make(map[uint32]*player),
			new_recv_channel(),
			make([]func([]byte,*Room),MAX_CMD_ID),
			make([]func(*Room),10),
		},
	}
}
func NilRoom()(*Room){
	return nil;
}
