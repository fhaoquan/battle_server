package old

import (
	"../../utils"
	"net"
	"fmt"
	"errors"
)

type UdpConnection struct {
	utils.ICachedData;
	Addr *net.UDPAddr
	conn *net.UDPConn;
}
func (me *UdpConnection)OnReturn(){
	if me.conn!=nil{
		me.conn.Close();
		me.conn=nil;
	}
}
func (me *UdpConnection)GetConn()*net.UDPConn{
	return me.conn;
}
func (me *UdpConnection)start()(err error){
	me.conn,err=net.ListenUDP("udp",me.Addr);
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
			s.conn=nil;
			start++;
			return s;
		}),
	};
	return m;
}

var(
	TheUDPConnManager=new_session_manager(utils.MaxRoomSize,utils.UdpListenStart);
)
