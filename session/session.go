package session

import (
	"net"
)
type S_session struct {
	V_user_id uint32;
	V_room_id uint32;
	V_conn net.Conn;
	v_udp_port uint16;
	F_send func(I_send_session_packet);
}
func (session *S_session)GetUserID()uint32{
	return session.V_user_id;
}
func (session *S_session)Start(
	conn net.Conn,
	finder func(uint32)(I_session_owner)){

	context:=session_proc_context{}
	context.start_session_proc(session, finder);
	session.F_send=func(pkt I_send_session_packet){
		context.send_pool<-pkt;
	}
}
func NewSession()(*S_session){
	return &S_session{
		0,
		0,
		nil,
		0,
		nil,
	}
}
func BuildSession(
	conn net.Conn,
	on_message func (rid uint32,pid uint32,len uint16,body []byte),
	get_msg_for_send func()){


}