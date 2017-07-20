package session

import (
	"net"
	"io"
	"../utils"
	"encoding/binary"
)

type kcp_session struct{
	uid uint32;
	rid uint32;
	con net.Conn;
}
func (me *kcp_session)ReadOne(packet *utils.Packet)(error){
	if _,e:=io.ReadFull(me.con,packet.BDY[0:2]);e!=nil{
		return e;
	}
	packet.LEN=binary.BigEndian.Uint16(packet.BDY[0:2])
	if _,e:=io.ReadFull(me.con,packet.BDY[2:packet.LEN]);e!=nil{
		return e;
	}
	return nil;
}
type kcp_session_proc struct {
	session *kcp_session;
}

type kcp_session_option struct {
	*kcp_session;
	e error;
}
func (me *kcp_session_option)Cast(f func(*kcp_session_option)(error))(*kcp_session_option){
	me.e=f(me);
	return me;
}
func (me *kcp_session_option)Loop(f func(*kcp_session_option)(error))(*kcp_session_option){
	for{
		me.e=f(me);
	}
	return me;
}

func empty_session(con net.Conn)(*kcp_session){
	return &kcp_session{
		uid:0,
		rid:0,
		con:con,
	}
}
func wait_room_id(session *kcp_session_option)(error){
	data:=make([]byte,1024);
	session.ReadOne(data);
	return nil;
}
func join_the_room(session *kcp_session_option)(error){
	return nil;
}
func wait_pkt_loop(session *kcp_session_option)(error){
	return nil;
}
func close_session(session *kcp_session_option)(error){
	return nil;
}
func tttt(){
	s:=&kcp_session_option{
		empty_session(nil),
		nil,
	};
	go func() {
		f:=wait_room_id;
		for{
			f=f(s);
		}
	}();
	utils.NewPacket(nil).Cast(func(packet *utils.Packet) {
		s.ReadOne(packet);
	}).Cast(func(packet *utils.Packet) {

	})
	go s.Cast(wait_room_id).
		Cast(join_the_room).
		Cast(wait_pkt_loop).
		Loop(func(session_option *kcp_session_option) error {
		return session_option.Cast(wait_pkt_loop).Cast(close_session).e;
	}).Cast(close_session);
	s.con.Close();
}
