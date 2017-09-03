package room

import (
	"../utils"
	"net"
	"fmt"
	"errors"
	"encoding/binary"
)

type UdpReq struct {
	utils.ICachedData;
	a net.Addr;
	l uint16;
	u uint32;
	r uint32;
	b []byte;
}
func (me *UdpReq)ReadAt(conn net.PacketConn)(err error){
	defer func(){
		if e:=recover();err!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}();
	_,a,e:=conn.ReadFrom(me.b);
	if e!=nil{
		return e;
	}
	me.l=binary.BigEndian.Uint16(me.b[4:6]);
	me.u=binary.BigEndian.Uint32(me.b[6:10]);
	me.r=binary.BigEndian.Uint32(me.b[10:14]);
	me.a=a;
	return nil;
}
func (me *UdpReq)GetALL()[]byte{
	return me.b;
}
func (me *UdpReq)GetLEN()uint16{
	return me.l;
}
func (me *UdpReq)GetUID()uint32{
	return me.u;
}
func (me *UdpReq)GetRID()uint32{
	return me.r;
}
func (me *UdpReq)GetRecvData()[]byte{
	return me.b[14:me.l+4];
}
func (me *UdpReq)GetAdr()net.Addr{
	return me.a;
}
