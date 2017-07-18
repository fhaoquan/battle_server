package package_glue
import (
	"../room"
	"../session"
)

type Room room.Room
func (c *Room)Join(s *session.S_session)bool{
	return (*room.Room)(c).Join((*Session)(s));
}
func (c *Room)OnMsg(packet *session.S_recv_session_packet){
	(*room.Room)(c).OnMsg(packet);
}