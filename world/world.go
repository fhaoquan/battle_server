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
func (w *World)AddNewRoom(r *room.Room)uint32{
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	w.room_id_seed++;
	r.SetID(w.room_id_seed);
	w.rooms[w.room_id_seed]=r;
	return w.room_id_seed;
}
func NewWorld()(*World){
	return &World{
		room_id_seed:0,
		rooms:make(map[uint32]*room.Room,1000),
	}
}