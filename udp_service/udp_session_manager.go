package udp_service

import (
	"../utils"
	"net"
	"fmt"
	"errors"
)

type UdpConnection struct {
	utils.ICachedData;
	Addr *net.UDPAddr
	*net.UDPConn;
}
func (me *UdpConnection)OnReturn(){
	if me.UDPConn!=nil{
		me.UDPConn.Close();
		me.UDPConn=nil;
	}
}
func (me *UdpConnection)start()(err error){
	me.UDPConn,err=net.ListenUDP("udp",me.Addr);
	if(err!=nil){
		return err;
	}
	return nil;
}

type session_manager struct {
	pool *utils.MemoryPool;
}
func (me *session_manager)Pop()(conn *UdpConnection,err error){
	t:=0;
	for{
		if(me.pool.Len()==0){
			return nil,errors.New("not free udp listener!");
		}
		c:=me.pool.Pop().(*UdpConnection);
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
			s:=new(UdpConnection);
			s.Addr,_=net.ResolveUDPAddr("udp",fmt.Sprint(":",start));
			s.ICachedData=impl;
			s.UDPConn=nil;
			start++;
			return s;
		}),
	};
	return m;
}

var(
	TheUDPConnManager=new_session_manager(utils.MaxRoomSize,utils.UdpListenStart);
)
