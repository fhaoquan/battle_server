package builder

import (
	"../sessions/kcp_session"
	"../world"
	"net"
)

func BuildKcpSession(conn net.Conn,world *world.World){
	kcp_session.NewSession(conn).StartAt(world);
}