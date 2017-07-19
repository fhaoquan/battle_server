package builder

import (
	"../session"
	"../world"
	"net"
)

func BuildKcpSession(conn net.Conn,world *world.World)(*session.S_session){
	s:=session.NewSession();
	s.Start(conn,func(rid uint32)session.I_session_owner{
		return &RoomGlue{world.FindRoom(rid)};
	});
	return s;
}