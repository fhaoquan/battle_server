package builder
import (
	"../room"
	"../session"
)

type RoomGlue struct{
	*room.Room
}
func (c *RoomGlue)Join(s *session.S_session)bool{
	return c.Room.Join((*SessionGlue)(s));
}
func (c *RoomGlue)OnMsg(packet *session.S_recv_session_packet){
	c.Room.OnMsg(packet)
	//(*room.Room)(c).OnMsg(packet);
}