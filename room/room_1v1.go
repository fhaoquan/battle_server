package room
import (
	"../utils"
	"../battle"
	"net"
	"github.com/sirupsen/logrus"
	"sync"
)
type kcp_connection_request struct {
	uid uint32;
	conn net.Conn;
}
type start_event struct{

}
type frame_event struct{
	frame uint32;
}
type BaseRoom struct {
	rid uint32;
	the_battle *battle.Battle;
	event_sig chan interface{};
	close_sig chan interface{};
	kcp_chan chan utils.IKcpRequest;
	udp_chan chan utils.IUdpRequest;
	once_start sync.Once;
	once_close sync.Once;
	wait sync.WaitGroup;
}
func (me *BaseRoom)GetID()(uint32){
	return me.rid;
}
func (me *BaseRoom)GetBattle()(*battle.Battle){
	return me.the_battle;
}
func new_base_room(the_battle *battle.Battle)(*BaseRoom){
	r:=&BaseRoom{
		0,
		the_battle,
		make(chan interface{},5),
		make(chan interface{},1),
		make(chan utils.IKcpRequest,16),
		make(chan utils.IUdpRequest,16),
		sync.Once{},
		sync.Once{},
		sync.WaitGroup{},
	}
	return r;
}
type room_player struct {
	uid uint32;
	name string;
	kcp_session *kcp_session;
	udp_session *udp_session;
}
type Room1v1 struct {
	*BaseRoom;
	p1 *room_player;
	p2 *room_player;
}
func (me *Room1v1)on_handler_result(rtn interface{}){
	switch rtn.(type){
	case nil:
		return ;
	case utils.IKcpResponse:
		me.on_kcp_response(rtn.(utils.IKcpResponse));
		rtn.(utils.IKcpResponse).Return();
	case []utils.IKcpResponse:
		for _,r:=range rtn.([]utils.IKcpResponse){
			me.on_kcp_response(r);
			r.Return();
		}
	case utils.IUdpResponse:
		me.on_udp_response(rtn.(utils.IUdpResponse));
		rtn.(utils.IUdpResponse).Return();
	case []utils.IUdpResponse:
		for _,r:=range rtn.([]utils.IUdpResponse){
			me.on_udp_response(r);
			r.Return();
		}
	case error:
		logrus.Error(rtn.(error));
	default:
		logrus.Error("unknown command response type! ");
		return ;
	}
}
func (me *Room1v1)on_packet(bdy []byte){
	switch bdy[0]{
	case 1:
		me.on_handler_result(me.the_battle.UpdateUnitMovement(bdy[1:]));
	case 2:
		me.on_handler_result(me.the_battle.UnitAttackStart(bdy[1:]));
	case 3:
		me.on_handler_result(me.the_battle.UnitAttackDone(bdy[1:]));
	case 4:
		me.on_handler_result(me.the_battle.CreateUnit(bdy[1:]));
	}
}
func (me *Room1v1)on_udp_response(r utils.IUdpResponse){
	if(r.IsBroadcast()){
		me.p1.udp_session.Send(r.GetSendData());
		me.p2.udp_session.Send(r.GetSendData());
	}else if me.p1.uid==r.GetUID(){
		me.p1.udp_session.Send(r.GetSendData());
	}else if me.p2.uid==r.GetUID(){
		me.p2.udp_session.Send(r.GetSendData());
	}
}
func (me *Room1v1)on_kcp_response(r utils.IKcpResponse){
	if(r.IsBroadcast()){
		me.p1.kcp_session.Send(r.GetSendData());
		me.p2.kcp_session.Send(r.GetSendData());
	}else if me.p1.uid==r.GetUID(){
		me.p1.kcp_session.Send(r.GetSendData());
	}else if me.p2.uid==r.GetUID(){
		me.p2.kcp_session.Send(r.GetSendData());
	}
}
func (me *Room1v1)on_kcp_message(r utils.IKcpRequest){
	me.on_packet(r.GetRecvData());
	r.Return();
}
func (me *Room1v1)on_udp_message(r utils.IUdpRequest){
	me.on_packet(r.GetRecvData());
	r.Return();
}
func (me *Room1v1)on_event(event interface{}){
	switch event.(type){
	case *kcp_connection_request:
		switch event.(*kcp_connection_request).uid {
		case me.p1.uid:
			me.p1.kcp_session=&kcp_session{event.(*kcp_connection_request).conn,me.p1.uid,}
			go me.kcp_recv_proc(me.p1.kcp_session);
		case me.p2.uid:
			me.p2.kcp_session=&kcp_session{event.(*kcp_connection_request).conn,me.p2.uid,}
			go me.kcp_recv_proc(me.p2.kcp_session);
		default:
			event.(*kcp_connection_request).conn.Close();
		}
	case *start_event:
		me.on_handler_result(me.the_battle.BroadcastBattleStart());
	case *frame_event:
		me.on_handler_result(me.the_battle.BroadcastBattleMovementData());
	}
}
func (me *Room1v1)OnKcpConnection(conn net.Conn,uid uint32){
	me.event_sig<-&kcp_connection_request{
		uid,conn,
	};
}
func (me *Room1v1)Start(){
	me.once_start.Do(func() {
		go me.start_proc();
	})
}
func (me *Room1v1)Close(why error){
	me.once_close.Do(func() {
		logrus.Error("room will closed for :",why)
		close(me.close_sig);
		me.wait.Wait();
	})
}
