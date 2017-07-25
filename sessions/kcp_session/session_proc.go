package kcp_session

import (
	"../../room"
	"../../utils"
	"../../world"
	"io"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"fmt"
	"github.com/pkg/errors"
	"net"
)
type empty_kcp_packet struct {
	pool *utils.PacketPool
}
func (me *empty_kcp_packet)Cast(f func(*empty_kcp_packet)interface{})interface{}{
	return f(me);
}
type kcp_packet struct {
	utils.I_cached_data;
	len uint16;
	uid uint32;
	rid uint32;
	buf []byte;
}
func (pkt *kcp_packet)recv(conn net.Conn)(interface{}){
	if _,e:=io.ReadFull(conn,pkt.buf[0:2]);e!=nil{
		return nil;
	}
	pkt.len=binary.BigEndian.Uint16(pkt.buf[0:2]);
	if _,e:=io.ReadFull(conn,pkt.buf[2:pkt.len]);e!=nil{
		return nil;
	}
	pkt.uid=binary.BigEndian.Uint32(pkt.buf[2:6]);
	pkt.rid=binary.BigEndian.Uint32(pkt.buf[6:10]);
	return pkt;
}
func (pkt *kcp_packet)withroom(f func(uint32)*room.Room)(interface{}){
	return f(pkt.rid);
}

type reader func ()(*kcp_packet,error);

func (rdr reader)find_room()*room.Room{
	rdr();
	return nil;
}
func build_reader(f func ()(*kcp_packet,error))reader{
	return f;
}



type command struct {
	kcp_packet;
	tar *room.Room;
}
func (c *command)to_room(r *room.Room){

}
func build_command(f func() *kcp_packet)(*command){

}
func read_one(pop func()*kcp_packet,con func()net.Conn)(interface{}){
	h:=pop();
	if _,e:=io.ReadFull(con(),h.buf[0:2]);e!=nil{
		return nil;
	}
	h.len=binary.BigEndian.Uint16(h.buf[0:2]);
	if _,e:=io.ReadFull(con(),h.buf[2:h.len]);e!=nil{
		return nil;
	}
	h.uid=binary.BigEndian.Uint32(h.buf[2:6]);
	h.rid=binary.BigEndian.Uint32(h.buf[6:10]);
	return h;
}
func handle_return(rtn interface{},checker func(r interface{})(int),handler ...func(r interface{})){
	handler[checker(rtn)](rtn);
}

func build_main_logic(s *session,w *world.World)(f func(interface{})interface{}){
	f=func(i interface{})interface{}{
		if i==nil{
			return i;
		}
		switch i.(type) {
		case *utils.PacketPool:
			f(i.(*utils.PacketPool).GetEmptyPkt().(*kcp_packet).recv(s.con));
		case *kcp_packet:
			f(i.(*kcp_packet).withroom(w.FindRoom));
		case *command:
			i.(*command).to_room(nil);
		case error:
			logrus.Debug(i);
		default:
			return nil;
		}
		return nil;
	};
	return ;
}
func ttt(s *session,w *world.World){
	pool:=utils.NewPacketPool(16, func(i utils.I_cached_data) utils.I_cached_data {
		return &kcp_packet{
			i,0,0,0,make([]byte,utils.MaxPktSize),
		}
	})
	on_err:=
		func(e error){
			logrus.Debug(e);
		}
	to_room:=
		func(command *command){
			switch r:=command.withroom(w.FindRoom);r.(type){
			case *command:
			case error:
				on_err(r.(error));
			}
		}
	to_command:=
		func(pkt *kcp_packet) {
			switch r:=pkt.withroom(w.FindRoom);r.(type){
			case *command:
				to_room(r.(*command))
			case error:
				on_err(r.(error));
			}
		}
	read_one:=
		func(){
			switch r:=pool.GetEmptyPkt().(*kcp_packet).recv(s.con); r.(type){
			case *kcp_packet:
				to_command(r.(*kcp_packet));
			case error:
				on_err(r.(error));
			}
		}
	start_as:=build_main_logic(s,w);
	for{
		start_as(pool)
	};

	f:=func(i interface{}){
		switch i.(type) {
		case *kcp_packet:
			func(i interface{}){
				switch i.(type) {
				case *command:
					i.(*command).to_room(nil);
				case error:
					func(i interface{}){
						logrus.Debug(i);
					}(i);
				}
			}(i.(*kcp_packet).withroom(w.FindRoom));
		case error:
			func(i interface{}){
				logrus.Debug(i);
			}(i);
		}
	};
	for{f(pool.GetEmptyPkt().(*kcp_packet).recv(s.con))};

	for{
		switch r:=pool.GetEmptyPkt().(*kcp_packet).recv(s.con); r.(type){
		case *kcp_packet:
			switch r=r.(*kcp_packet).withroom(w.FindRoom);r.(type){
			case *command:
				switch r=r.(*kcp_packet).withroom(w.FindRoom);r.(type){
				case *kcp_packet:
				case error:
					logrus.Error(r);
				}
			case error:
				logrus.Error(r);
			}
		case error:
			logrus.Error(r);
		}
	}

}