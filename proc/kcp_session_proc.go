package proc

import (
	"net"
	"../utils"
	"../room"
	"../world"
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"fmt"
	"io"
)
type i_kcp_packet interface {
	ReadAt(net.Conn)(error);
	GetUID()uint32;
	GetRID()uint32;
	GetTID()uint8;
	GetBDY()[]byte;
}
type kcp_session_context struct {
	head32 []byte;
	con net.Conn;
	the_player *room.Player;
	the_world *world.World;
	the_room *room.Room;
	pool *utils.MemoryPool;
}
func (s *kcp_session_context)send_proc(sig chan interface{})(error){
	c:=make(chan []byte,16);
	s.the_player.SetKcpSender(c);
	for{
		select {
		case pkt:=<-c:
			s.con.Write(pkt);
		case <-sig:
			return nil;
		}
	}
}
func (s *kcp_session_context)main_proc(msg chan utils.IDataOwner,sig chan interface{})(error){
	for{
		select {
		case cache:=<-msg:
			s.the_room.OnKCP(
				0,
				cache.GetUserData().(i_kcp_packet).GetUID(),
				cache.GetUserData().(i_kcp_packet).GetRID(),
				cache.GetUserData().(i_kcp_packet).GetBDY(),
			);
			cache.Return();
		case <-sig:
			return nil;
		}
	}
}
func (s *kcp_session_context)recv_proc(msg chan utils.IDataOwner,sig chan interface{})(error){
	for{
		b:=s.pool.PopOne();
		if e:=s.read(b.GetUserData().(*kcp_packet));e!=nil{
			b.Return();
			return e;
		}
		select {
		case msg<-b:
		case sig:
			b.Return();
			return nil;
		}
	}
}
func (s *kcp_session_context)handle_first_packet(packet *kcp_packet)(error){
	if e:=s.read(packet);e!=nil{
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
func (s *kcp_session_context)StartAt(w *world.World){
	s.the_world=w;
	go func(){
		b:=s.pool.PopOne();
		if e:=s.handle_first_packet(b.GetUserData().(*kcp_packet));e!=nil{
			b.Return();
			s.con.Close();
		}
		msg:=make(chan utils.IDataOwner,16);
		msg<-b;
		sig:=make(chan interface{},3);
		go func(err error){
			if(err!=nil){
				logrus.Error(err);
			}
		}(s.send_proc(sig));
		go func(err error){
			if(err!=nil){
				logrus.Error(err);
			}
		}(s.main_proc(msg,sig));
		go func(err error){
			sig<-1;
			sig<-1;
			s.con.Close();
			if(err!=nil){
				logrus.Error(err);
			}
		}(s.recv_proc(msg,sig));
	}();
}
func StartKcpSession(con net.Conn,wld world.World){

}
