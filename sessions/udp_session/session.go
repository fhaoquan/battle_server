package udp_session

import (
	"net"
	"../../utils"
	"../../room"
	"errors"
	"github.com/sirupsen/logrus"
	"encoding/binary"
	"sync"
)
type Session struct {
	pool *utils.MemoryPool;
	conn *udp_connection;
	the_room *room.Room;
	req chan utils.IUdpRequest;
	res chan utils.IUdpResponse;
	once sync.Once;
	wait sync.WaitGroup;
}
func (s *Session)GetAddr()*net.UDPAddr{
	return s.conn.addr;
}
func (s *Session)send_proc(c <-chan utils.IUdpResponse)(error){
	s.the_room.ForEachPlayer(func(player *room.Player)bool{
		player.SetUDPSender(c);
		return true;
	});
	for{
		select {
		case pkt,ok:=<-c:
			if !ok{
				pkt.Return();
				return nil;
			}
			if _,e:=s.conn.WriteTo(pkt.GetBDY(),pkt.GetAdr());e!=nil{
				pkt.Return();
				return e;
			}
			pkt.Return();
		}
	}
}
func (s *Session)main_proc(cmd <-chan utils.IUdpRequest)(error){
	for{
		select {
		case pkt,ok:=<-cmd:
			if !ok{
				pkt.Return();
				return nil;
			}else{
				s.the_room.OnUDP(
					pkt.GetAdr(),
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
func (s *Session)read(packet *udp_packet)(error){
	_,a,e:=s.conn.ReadFrom(packet.b);
	if e!=nil{
		return e;
	}
	if(binary.BigEndian.Uint32(packet.b[0:4])!=123454321){
		return e;
	}
	packet.a=a;
	packet.l=binary.BigEndian.Uint16(packet.b[4:6]);
	packet.u=binary.BigEndian.Uint32(packet.b[6:10]);
	packet.r=binary.BigEndian.Uint32(packet.b[10:14]);
	return nil;
}
func (s *Session)recv_proc(c chan utils.IUdpRequest)(error){
	f:= func(p utils.IUdpRequest)(ok bool){
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
		p:=s.pool.Pop().(utils.IUdpRequest);
		switch e:=p.ReadAt(s.conn);e.(type){
		case nil:
			if !f(p){
				return nil;
			}
		case net.Error:
			p.Return();
			if e.(net.Error).Temporary(){
				if t++;t>=5 {
					return errors.New("5 times temporary error!");
				}else{
					continue;
				}
			}else{
				return e;
			}
		default:
			p.Return();
			return e;
		}
		t=0;
	}
}
func (s *Session)handle_err(e error){
	logrus.Error(e);
}
func (s *Session)StartAt(room *room.Room){
	s.req=make(chan utils.IUdpRequest,8);
	s.res=make(chan utils.IUdpResponse,8);
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
}
func (s *Session)Close(){
	s.once.Do(func() {
		close(s.req);
		close(s.res);
		s.wait.Wait();
		s.conn.Return();
	})
}
func NewSession()(*Session,error){
	if c,e:=the_session_manager.pop();e!=nil{
		return nil,e;
	}else{
		return &Session{conn:c},nil;
	}
}
