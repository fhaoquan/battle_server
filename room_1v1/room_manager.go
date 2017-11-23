package room_1v1

import (
	"sync"
)

var rooms				=new(sync.Map);
var room_id_seed		=uint32(1000);

type RoomBuildContext struct {
	Lifecycle		int;
	SuddenDeath		int;
	WinScore		int;
}

func add_room(k uint32,r *room)  {
	rooms.Store(k,r);
}
func del_room(k uint32){
	rooms.Delete(k);
}
func get_room(k uint32)(*room){
	r,y:=rooms.Load(k);
	if y{
		return r.(*room);
	}else{
		return nil;
	}
}
func CountRoom()(int){
	return 1000;
}
