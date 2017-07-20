package room

import (
	"../utils"
	"time"
)

const MAX_CMD_ID  = 255;
type i_session interface {
	GetUserID()uint32;
}
type i_battle interface {

}
type Room struct{
	id uint32;
	battle i_battle;
	players map[uint32]*player;
	recv_msg_chan chan room_message;
	cmd_handlers []func([]byte,*Room);
	timer_handlers []func(*Room);
}
type room_message interface {
	utils.I_cached_data;
	GetMsgBody()[]byte;
	GetUser()uint32;
	GetRoom()uint32;
}
func (r *Room)SetID(v uint32){
	r.id=v;
}
func (r *Room)OnPkt(uid uint32,rid uint32,body []byte){

}
func (r *Room)OnMsg(msg room_message){
	r.recv_msg_chan<-msg;
}
func (r *Room)Join(s i_session)bool{
	if player,ok:=r.players[s.GetUserID()];ok{
		player.session=s;
		return true;
	}
	return false;
}
func (r *Room)Broadcast(data []byte,len int){

}
func (r *Room)Start(){
	timer:=make(chan int,1);
	go func(){
		for{
			select {
			case msg:=<-r.recv_msg_chan:
				if(r.cmd_handlers[msg.GetMsgBody()[0]]!=nil){
					r.cmd_handlers[msg.GetMsgBody()[0]](msg.GetMsgBody()[1:],r);
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
			nil,
			make(map[uint32]*player),
			make(chan room_message,128),
			make([]func([]byte,*Room),MAX_CMD_ID),
			make([]func(*Room),10),
		},
	}
}
func NilRoom()(*Room){
	return nil;
}
