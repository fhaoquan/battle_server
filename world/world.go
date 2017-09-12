package world

import (
	"sync"
	"../room"
	"../utils"
)

type World struct {
	m			sync.RWMutex;
	rooms		map[uint32]*room.Room1v1;
	id_seed		uint32
	pool		*utils.MemoryPool;
}
func (w *World)CountRoom()(int){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	return len(w.rooms);
}
func (w *World)ForEachRoom(f func(*room.Room1v1)bool){
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
func (w *World)FindRoom(rid uint32)(*room.Room1v1){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	if room,ok:=w.rooms[rid];ok{
		return room;
	}
	return nil;
}
func (w *World)AddNewRoom(r *room.Room1v1){
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	r.SetID(w.id_seed);
	w.rooms[r.GetID()]=r;
	w.id_seed++;
	r.Start(w);
}
func (w *World)DelRoom(r *room.Room1v1)  {
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	delete(w.rooms,r.GetID());
}
func NewWorld()(*World){
	return &World{
		rooms:make(map[uint32]*room.Room1v1,1000),
		id_seed:10000,
		pool:utils.NewMemoryPool(64, func(impl utils.ICachedData) utils.ICachedData {
			return &utils.KcpReq{
				impl,make([]byte,utils.MaxPktSize),
			}
		}),
	}
}