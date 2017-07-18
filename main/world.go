package main

import (
	"sync"
	"room"
)

type world struct {
	m sync.RWMutex;
	room_id_seed uint32;
	rooms map[uint32]*room.Room;
}
func (w *world)FindRoom(rid uint32)(*room.Room){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	if room,ok:=w.rooms[rid];ok{
		return room;
	}
	return nil;
}
func (w *world)AddNewRoom(b *room.S_room_builder)uint32{
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	w.room_id_seed++;
	b.ID(w.room_id_seed);
	w.rooms[w.room_id_seed]=b.Build();
	return w.room_id_seed;
}
func new_world()(*world){
	return &world{
		room_id_seed:0,
		rooms:make(map[uint32]*room.Room,1000),
	}
}