package session

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/pkg/errors"
	"encoding/binary"
	"io"
	"../utils"
)
type I_send_session_packet interface {
	utils.I_cached_data;
	GetMsgBody()[]byte;
}
type S_recv_session_packet struct {
	utils.I_cached_data;
	len uint16;
	uid uint32;
	rid uint32;
	data []byte;
}
func (pkt *S_recv_session_packet)Clear(){
	pkt.len=0;
	pkt.uid=0;
	pkt.rid=0;
}
func (pkt *S_recv_session_packet)GetRoom()uint32{
	return pkt.rid;
}
func (pkt *S_recv_session_packet)GetUser()uint32{
	return pkt.uid;
}
func (pkt *S_recv_session_packet)GetMsgBody()[]byte{
	return pkt.data[11:pkt.len];
}
type I_session_owner interface {
	OnMsg(*S_recv_session_packet);
	Join(*S_session)bool;
	OnPkt(uint32,uint32,[]byte);
}
type session_proc_context struct{
	send_pool chan I_send_session_packet;
	send_ctrl chan interface{};
	recv_ctrl chan interface{};
	base_room I_session_owner;
}
func (context *session_proc_context)send_proc(session *S_session){
	data:=make([]byte,utils.MaxPktSize);
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
			close(context.send_pool);
		}();

		for{
			select {
			case msg,ok:=<-context.send_pool:
				if(!ok){
					return errors.New("send channel closed");
				}
				l:=uint16(len(msg.GetMsgBody()));
				binary.BigEndian.PutUint16(data[0:2],l);
				_,err:=session.V_conn.Write(append(data[0:2],msg.GetMsgBody()...)[0:l]);
				msg.ReturnToPool();
				if(err!=nil){
					return err;
				}
			case <-context.send_ctrl:
				return errors.New("send proc stop by command");
			}
		}

	}();
	if(err!=nil){
		logrus.Error(err);
	}
}
func (context* session_proc_context)recv_proc(
	session *S_session,
	finder func(uint32)(I_session_owner),
){
	pool:=utils.NewPacketPool(64,func(icd utils.I_cached_data)utils.I_cached_data{
		return &S_recv_session_packet{
			icd,
			0,
			0,
			0,
			make([]byte,utils.MaxPktSize),
		};
	});
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		for{
			d:=pool.GetEmptyPkt().(*S_recv_session_packet);
			if _,err:=io.ReadFull(session.V_conn,d.data[0:2]);err!=nil{
				return err;
			}
			l:=binary.BigEndian.Uint16(d.data[0:2]);
			if _,err:=io.ReadFull(session.V_conn,d.data[2:l+2]);err!=nil{
				return  err;
			}
			if(session.V_user_id!=0){
				session.V_user_id=binary.BigEndian.Uint32(d.data[3:4+3]);
			}
			if(session.V_room_id!=0){
				session.V_room_id=binary.BigEndian.Uint32(d.data[7:4+7]);
			}
			if(context.base_room==nil){
				r:=finder(session.V_room_id)
				if(r==nil){
					return errors.New("can't find room");
				}
				if(r.Join(session)){
					return errors.New("can't find room");
				}
				context.base_room=r;
			}
			switch(d.data[11]){
			case 1:
				session.v_udp_port=binary.BigEndian.Uint16(d.data[12:14]);
			default:
				context.base_room.OnMsg(d);
			}
		}
	}();
	if(err!=nil){
		logrus.Error(err);
	}
}
func (context* session_proc_context)start_session_proc(
	session *S_session,
	finder func(uint32)(I_session_owner)){
	context.recv_ctrl=make(chan interface{},1);
	context.send_ctrl=make(chan interface{},1);
	context.send_pool=make(chan I_send_session_packet,32);
	go context.recv_proc(session,finder);
	go context.send_proc(session);
}
