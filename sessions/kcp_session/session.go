package kcp_session

import "net"

type session struct {
	uid uint32;
	rid uint32;
	con net.Conn;
}

func NewSession(con net.Conn)(s *session){
	return &session{
		con:con,
	}
}
