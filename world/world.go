package world

import (
	"sync"
	"../room"
	"net"
)

type World struct {
	m sync.RWMutex;
	rooms map[uint32]*room.Room;
}
func (w *World)CountRoom()(int){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	return len(w.rooms);
}
func (w *World)ForEachRoom(f func(*room.Room)bool){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	for _,v:=range w.rooms{
		if !f(v){
			return ;
		}
	}
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
func (w *World)AddNewRoom(r *room.Room){
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	w.rooms[r.GetID()]=r;
}
func (w *World)OnNewKCPConnection(conn net.Conn){
}
func NewWorld()(*World){
	return &World{
		rooms:make(map[uint32]*room.Room,1000),
	}
}