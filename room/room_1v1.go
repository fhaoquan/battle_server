package room
import (
	"../utils"
	"../battle"
	"../server/kcp_server"
	"github.com/sirupsen/logrus"
	"sync"
	"net"
	"runtime/debug"
	"fmt"
	"errors"
)
type kcp_connection_request struct {
	session		*kcp_server.KcpSession;
	uid			uint32;
}
type kcp_session_closed struct{
	p 			*room_player;
	s 			*kcp_server.KcpSession;
}
type I_RoomManager interface{
	AddNewRoom(*Room1v1)
	DelRoom(*Room1v1)
}
type start_event struct{
	state int;
}
type frame_event struct{
	frame uint32;
}
type BaseRoom struct {
	rid			uint32;
	the_battle	*battle.Battle;
	event_sig	chan interface{};
	close_sig	chan interface{};
	kcp_chan	chan *utils.KcpReq;
	udp_chan	chan *utils.UdpReq;
	udp_sender	*net.UDPConn;
	once_start	sync.Once;
	once_close	sync.Once;
	wait		sync.WaitGroup;
	manager		I_RoomManager;
}
func (me *BaseRoom)SetID(v uint32){
	me.rid=v;
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
		make(chan *utils.KcpReq,16),
		make(chan *utils.UdpReq,16),
		nil,
		sync.Once{},
		sync.Once{},
		sync.WaitGroup{},
		nil,
	}
	r.udp_sender,_=net.ListenUDP("udp",&net.UDPAddr{net.IPv4zero,0,""});
	return r;
}
type room_player struct {
	uid				uint32;
	name			string;
	kcp_session		*kcp_server.KcpSession;
	peer_udp_addr	net.Addr;
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
	case *utils.KcpRes:
		me.on_kcp_response(rtn.(*utils.KcpRes));
		rtn.(*utils.KcpRes).Return();
	case []*utils.KcpRes:
		for _,r:=range rtn.([]*utils.KcpRes){
			me.on_kcp_response(r);
			r.Return();
		}
	case *utils.UdpRes:
		me.on_udp_response(rtn.(*utils.UdpRes));
		rtn.(*utils.UdpRes).Return();
	case []*utils.UdpRes:
		for _,r:=range rtn.([]*utils.UdpRes){
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
func (me *Room1v1)on_packet(who uint32,bdy []byte){
	switch bdy[0]{
	case utils.CMD_pingpong:
		me.on_handler_result(me.the_battle.Pong(who,bdy[1:]));
	case utils.CMD_unit_movment:
		me.on_handler_result(me.the_battle.UpdateUnitMovement(bdy[1:]));
	case utils.CMD_attack_start:
		me.on_handler_result(me.the_battle.UnitAttackStart(bdy[1:]));
	case utils.CMD_attack_done:
		me.on_handler_result(me.the_battle.UnitAttackDone(bdy[1:]));
	case utils.CMD_create_unit:
		me.on_handler_result(me.the_battle.CreateUnit(who,bdy[1:]));
	case utils.CMD_unit_destory:
		me.on_handler_result(me.the_battle.UnitDestory(bdy[1:]));
	}
}
func (me *Room1v1)on_udp_response(r *utils.UdpRes){
	defer func() {
		r.Return();
		if e:=recover();e!=nil{
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	if(r.Broadcast){
		if me.p1.peer_udp_addr!=nil{
			me.udp_sender.WriteTo(r.GetSendData(),me.p1.peer_udp_addr);
		}
		if me.p2.peer_udp_addr!=nil{
			me.udp_sender.WriteTo(r.GetSendData(),me.p2.peer_udp_addr);
		}
	}else if me.p1.uid==r.GetUID(){
		if me.p1.peer_udp_addr!=nil{
			me.udp_sender.WriteTo(r.GetSendData(),me.p1.peer_udp_addr);
		}
	}else if me.p2.uid==r.GetUID(){
		if me.p2.peer_udp_addr!=nil{
			me.udp_sender.WriteTo(r.GetSendData(),me.p2.peer_udp_addr);
		}
	}
}
func (me *Room1v1)on_kcp_response(r *utils.KcpRes){
	defer func() {
		r.Return();
		if e:=recover();e!=nil{
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	if(r.IsBroadcast()){
		if me.p1.kcp_session!=nil{
			me.p1.kcp_session.Send(r.GetSendData());
		}
		if me.p2.kcp_session!=nil{
			me.p2.kcp_session.Send(r.GetSendData());
		}
	}else if me.p1.uid==r.GetUID(){
		if me.p1.kcp_session!=nil{
			me.p1.kcp_session.Send(r.GetSendData());
		}
	}else if me.p2.uid==r.GetUID(){
		if me.p2.kcp_session!=nil{
			me.p2.kcp_session.Send(r.GetSendData());
		}
	}
}
func (me *Room1v1)on_kcp_message(r *utils.KcpReq){
	defer r.Return();
	me.on_packet(r.GetUID(),r.GetMsgBody());
}
func (me *Room1v1)on_udp_message(r *utils.UdpReq){
	defer r.Return();
	if me.p1.uid==r.GetUID(){
		me.p1.peer_udp_addr=r.GetAdr();
	}else if me.p2.uid==r.GetUID(){
		me.p2.peer_udp_addr=r.GetAdr();
	}else{
		return ;
	}
	me.on_packet(r.GetUID(),r.GetMsgBody());
}
func (me *Room1v1)on_event(event interface{}){
	switch event.(type){
	case *kcp_connection_request:
		switch event.(*kcp_connection_request).uid {
		case me.p1.uid:
			me.p1.kcp_session=event.(*kcp_connection_request).session;
			go me.room_kcp_recv_proc(me.p1);
		case me.p2.uid:
			me.p2.kcp_session=event.(*kcp_connection_request).session;
			go me.room_kcp_recv_proc(me.p2);
		default:
			event.(*kcp_connection_request).session.Close(false);
		}
	case *kcp_session_closed:
		e:=event.(*kcp_session_closed);
		if e.p.kcp_session==e.s{
			e.p.kcp_session=nil;
		}
	case *start_event:
		switch(event.(*start_event).state){
		case 0:
			me.on_handler_result(me.the_battle.BroadcastBattleWaitingStart());
		case 1:
			me.on_handler_result(me.the_battle.BroadcastBattleStart());
			me.on_handler_result(me.the_battle.BroadcastBattleAll());
		}
	case *frame_event:
		switch res:=me.the_battle.CheckBattleEnd();res.(type) {
		case nil:
			me.on_handler_result(me.the_battle.BroadcastBattleMovementData(me.p1.uid,me.p2.uid));
			me.on_handler_result(me.the_battle.BroadcastBattleMovementData(me.p2.uid,me.p1.uid));
		default:
			me.on_handler_result(res);
			me.Close(errors.New("the battle complated!"));
		}
	}
}
func (me *Room1v1)OnKcpSession(uid uint32,session *kcp_server.KcpSession){
	me.event_sig<-&kcp_connection_request{session,uid};
}

func (me *Room1v1)Start(manager I_RoomManager){
	me.manager=manager;
	go me.once_start.Do(func() {
		me.start_proc();
	})
}
func (me *Room1v1)Close(why error){
	go me.once_close.Do(func() {
		logrus.Error("room ",me.rid," will closed for :",why);
		close(me.close_sig);
		me.wait.Wait();
		if(me.manager!=nil){
			me.manager.DelRoom(me);
		}
		logrus.Error("room ",me.rid," closed ");
	})
}
