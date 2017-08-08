package command

import (
	"../room"
)

type Commamd struct {
	base_room *room.Room;
}
func (cmd *Commamd)SetRoom(r *room.Room){
	cmd.base_room=r;
}
