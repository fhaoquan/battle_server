package room

import (
	"../utils"
	"time"
	"net"
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
func (r *Room)SetID(v uint32){
	r.id=v;
}
func (r *Room)KCPSend(uid uint32,data []byte,len int){
}
func (r *Room)KCPBroadcast(data []byte,len int){
}
func (r *Room)UDPSend(uid uint32,data []byte,len int){
}
func (r *Room)UDPBroadcast(data []byte,len int){
}
func (r *Room)Start(){
	if(r.started){
		return ;
	}
	r.started=true;
	timer:=make(chan int,1);
	go func(){
		for{
			select {
			case f:=<-r.kcp_chan:
				f(func(uid uint32, rid uint32, bdy []byte)bool{
					if(r.cmd_handlers[bdy[0]]!=nil){
						r.cmd_handlers[bdy[0]](bdy[1:],r);
					}
					return true;
				})
			case f:=<-r.udp_chan:
				f(func(adr net.Addr, uid uint32, rid uint32, bdy []byte)bool{
					if(r.cmd_handlers[bdy[0]]!=nil){
						r.cmd_handlers[bdy[0]](bdy[1:],r);
					}
					return true;
				})
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
}
func NewRoom()(*S_room_builder){
	return &S_room_builder{
		&Room{
			0,
			false,
			make(map[uint32]*player),
			new_recv_channel(),
			make(chan *kcp_message,128),
			make([]func([]byte,*Room),MAX_CMD_ID),
			make([]func(*Room),10),
		},
	}
}
func NilRoom()(*Room){
	return nil;
}
