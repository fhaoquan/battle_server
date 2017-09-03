package world

import (
	"sync"
	"../room"
	"net"
	"io"
	"encoding/binary"
)

type World struct {
	m sync.RWMutex;
	rooms map[uint32]*room.BattleRoom1v1;
}
func (w *World)CountRoom()(int){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	return len(w.rooms);
}
func (w *World)ForEachRoom(f func(*room.BattleRoom1v1)bool){
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
func (w *World)FindRoom(rid uint32)(*room.BattleRoom1v1){
	defer func(){
		w.m.RUnlock();
	}();
	w.m.RLock();
	if room,ok:=w.rooms[rid];ok{
		return room;
	}
	return nil;
}
func (w *World)AddNewRoom(r *room.BattleRoom1v1){
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	w.rooms[r.GetID()]=r;
}
func (w *World)OnNewKCPConnection(conn net.Conn){
	go func(){
		b:=make([]byte,16);
		io.ReadFull(conn,b);
		r:=binary.BigEndian.Uint32(b[0:4]);
		u:=binary.BigEndian.Uint32(b[4:8]);
		w.FindRoom(r).OnKcpConnection(conn,u);
	}();
}
func NewWorld()(*World){
	return &World{
		rooms:make(map[uint32]*room.BattleRoom1v1,1000),
	}
}