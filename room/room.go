package room

import (
	"../utils"
	"time"
	"net"
)

const MAX_CMD_ID  = 255;
type i_session interface {
	GetUserID()uint32;
}


type Room struct{
	id uint32;
	players map[uint32]*player;
	*recv_channel;
	send_msg_chan chan *room_message;
	cmd_handlers []func([]byte,*Room);
	timer_handlers []func(*Room);
}

func (r *Room)SetID(v uint32){
	r.id=v;
}
func (r *Room)Join(s i_session)bool{
	if player,ok:=r.players[s.GetUserID()];ok{
		player.session=s;
		return true;
	}
	return false;
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
	timer:=make(chan int,1);
	go func(){
		for{
			select {
			case msg:=<-r.recv_msg_chan:
				if(r.cmd_handlers[msg.bdy[0]]!=nil){
					r.cmd_handlers[msg.bdy[0]](msg.bdy[1:],r);
				}
				msg.ReturnToPool();
				break;
			case tid:=<-timer:
				if(r.timer_handlers[tid]!=nil){
					r.timer_handlers[tid](r);
				}
				break;
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
			make(map[uint32]*player),
			new_recv_channel(),
			make(chan *room_message,128),
			make([]func([]byte,*Room),MAX_CMD_ID),
			make([]func(*Room),10),
		},
	}
}
func NilRoom()(*Room){
	return nil;
}
