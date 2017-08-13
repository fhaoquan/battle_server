package udp_session

import (
	"../../utils"
	"net"
	"fmt"
	"errors"
)

type udp_connection struct {
	addr *net.UDPAddr
	utils.ICachedData;
	net.PacketConn;
}
func (me *udp_connection)OnReturn(){
	if me.PacketConn!=nil{
		me.PacketConn.Close();
		me.PacketConn=nil;
	}
}
func (me *udp_connection)start()(err error){
	me.PacketConn,err=net.ListenUDP("udp",me.addr);
	if(err!=nil){
		return err;
	}
	return nil;
}

type session_manager struct {
	pool *utils.MemoryPool;
}
func (me *session_manager)pop()(conn *udp_connection,err error){
	t:=0;
	for{
		if(me.pool.Len()==0){
			return nil,errors.New("not free udp listener!");
		}
		c:=me.pool.Pop().(*udp_connection);
		if e:=c.start();e==nil{
			return c,nil;
		}
		if t++;t<6{
			c.Return();
		}else {
			return nil,errors.New("tomany error times at listen udp");
		}
	}
}

func new_session_manager(size int,start int)(*session_manager){
	m:=&session_manager{
		pool:utils.NewMemoryPool(size, func(impl utils.ICachedData)utils.ICachedData{
			s:=new(udp_connection);
			s.addr,_=net.ResolveUDPAddr("udp",fmt.Sprint(":",start));
			s.ICachedData=impl;
			s.PacketConn=nil;
			start++;
			return s;
		}),
	};
	return m;
}

var(
	the_session_manager=new_session_manager(utils.MaxRoomSize,utils.UdpListenStart);
)
