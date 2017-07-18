package room

import (
	"sync"
)

type S_Hall struct {
	m sync.RWMutex;
	room_id_seed uint32;
	rooms map[uint32]*Room;
}
func (context *S_Hall)FindRoom(rid uint32)(*Room){
	defer func(){
		context.m.RUnlock();
	}();
	context.m.RLock();
	if room,ok:=context.rooms[rid];ok{
		return room;
	}
	return nil;
}

func (context *S_Hall)AddNewRoom(r *Room)uint32{
	defer func(){
		context.m.Unlock();
	}();
	context.m.Lock();
	context.room_id_seed++;
	r.id=context.room_id_seed;
	context.rooms[r.id]=r;
	return r.id;
}

func NewHall()(*S_Hall){
	return &S_Hall{
		rooms:make(map[uint32]*Room),
		room_id_seed:0,
	}
}
