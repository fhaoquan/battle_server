package world

import (
	"sync"
	"../room"
	"../utils"
	"net"
	"github.com/sirupsen/logrus"
)

type World struct {
	m sync.RWMutex;
	rooms map[uint32]*room.Room1v1;
	pool *utils.MemoryPool;
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
	w.rooms[r.GetID()]=r;
	r.Start(w);
}
func (w *World)DelRoom(r *room.Room1v1)  {
	defer func(){
		w.m.Unlock();
	}();
	w.m.Lock();
	delete(w.rooms,r.GetID());
}
func (w *World)OnNewKCPConnection(conn net.Conn){
	go func(){
		req:=w.pool.Pop().(*utils.KcpReq);
		if e:=req.ReadAt(conn);e!=nil{
			logrus.Error("on kcp new connection",e);
			req.Return();
			conn.Close();
			return ;
		}
		if r:=w.FindRoom(req.RID);r!=nil{
			r.OnKcpConnection(conn,req);
		}else{
			logrus.Error("on kcp new connection cant find room=",req.RID);
			req.Return();
			conn.Close();
			return ;
		}
	}();
}
func NewWorld()(*World){
	return &World{
		rooms:make(map[uint32]*room.Room1v1,1000),
		pool:utils.NewMemoryPool(64, func(impl utils.ICachedData) utils.ICachedData {
			return &utils.KcpReq{
				impl,0,0,0,make([]byte,utils.MaxPktSize),
			}
		}),
	}
}