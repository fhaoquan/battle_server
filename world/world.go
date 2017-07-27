package world

import (
	"sync"
	"../room"
)

type World struct {
	m sync.RWMutex;
	room_id_seed uint32;
	rooms map[uint32]*room.Room;
}
func (w *World)FindRoom(rid uint32)(*room.Room){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	if room,ok:=w.rooms[rid];ok{
		return room;
	}
	return nil;
}
func (w *World)AddNewRoom(new_room func(id uint32)*room.Room){
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	r:=new_room(w.room_id_seed);
	w.rooms[r.GetID()]=r;
	w.room_id_seed=r.GetID()+1;
}
func NewWorld()(*World){
	return &World{
		room_id_seed:0,
		rooms:make(map[uint32]*room.Room,1000),
	}
}