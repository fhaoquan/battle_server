package kcp_session

import (
	"net"
	"errors"
	"../../world"
	"../../room"
	"../../utils"
	"../packet"
	"github.com/sirupsen/logrus"
	"sync"
)

type Session struct {
	head32 []byte;
	con net.Conn;
	the_player *room.Player;
	the_world *world.World;
	the_room *room.Room;
	pool *utils.MemoryPool;
	req chan packet.IKcpRequest;
	res chan packet.IKcpResponse;
	once sync.Once;
	wait sync.WaitGroup;
}

func (s *Session)send_proc(cmd <-chan packet.IKcpResponse)(error){
	for{
		select {
		case pkt,ok:=<-cmd:
			if !ok{
				pkt.Return();
				return nil;
			}
			if _,e:=s.con.Write(pkt.GetBDY());e!=nil{
				pkt.Return();
				return e;
			}
			pkt.Return();
		}
	}
}
func (s *Session)main_proc(cmd <-chan packet.IKcpRequest)(error){
	for{
		select {
		case pkt,ok:=<-cmd:
			if !ok{
				pkt.Return();
				return nil;
			}else{
				s.the_room.OnKCP(
					pkt.GetLEN(),
					pkt.GetUID(),
					pkt.GetRID(),
					pkt.GetBDY(),
				);
				pkt.Return();
			}
		}
	}
}
func (s *Session)recv_proc(c chan<- packet.IKcpRequest)(error){
	f:= func(p packet.IKcpRequest)(ok bool){
		defer func(){
			if recover()!=nil{
				ok=false;
			}
		}();
		c<-p;
		return true;
	}
	t:=0;
	for{
		r:=s.pool.Pop().(*KcpReq);
		switch e:=r.ReadAt(s.con);e.(type){
		case nil:
			if !f(r){
				return nil;
			}
		case net.Error:
			if e.(net.Error).Temporary(){
				r.Return();
				if t++;t>=5 {
					return errors.New("5 times temporary error!");
				}else{
					continue;
				}
			}else{
				r.Return();
				return e;
			}
		default:
			r.Return();
			return e;
		}
		t=0;
	}
}
func (s *Session)handle_first_packet(packet *KcpReq)(error){
	if e:=packet.ReadAt(s.con);e!=nil{
		return e;
	}
	s.the_room=s.the_world.FindRoom(packet.r);
	s.the_room.ForOnePlayer(packet.u, func(player *room.Player){
		s.the_player=player;
	})
	if s.the_room!=nil&&s.the_player!=nil{
		return errors.New("can not find room");
	}
	return nil;
}
func (s *Session)StartAt(w *world.World){
	go func() {
		s.the_world=w;
		s.req=make(chan packet.IKcpRequest,8);
		s.res=make(chan packet.IKcpResponse,8);
		b := s.pool.Pop().(*KcpReq);
		if e := s.handle_first_packet(b); e != nil {
			b.Return();
			s.Close();
		}
		s.req<-b;
		go func(){
			s.wait.Add(1);
			if e:=s.send_proc(s.res);e!=nil{
				logrus.Error(e);
			}
			s.wait.Done();
			s.Close();
		}();
		go func(){
			s.wait.Add(1);
			if e:=s.main_proc(s.req);e!=nil{
				logrus.Error(e);
			}
			s.wait.Done();
			s.Close();
		}();
		go func(){
			s.wait.Add(1);
			if e:=s.recv_proc(s.req);e!=nil{
				logrus.Error(e);
			}
			s.wait.Done();
			s.Close();
		}();
	}();
}
func (s *Session)Close(){
	s.once.Do(func() {
		close(s.req);
		close(s.res);
		s.wait.Wait();
		s.con.Close();
	})
}

func NewSession(con net.Conn)(s *Session){
	return &Session{
		con:con,
		pool:utils.NewMemoryPool(16, func(impl utils.ICachedData) utils.ICachedData {
			return &KcpReq{
				impl,0,0,0,make([]byte,utils.MaxPktSize),
			}
		}),
	}
}
