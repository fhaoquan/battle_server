package utils

import (
	"io"
	"net"
	"fmt"
	"errors"
	"encoding/binary"
	"time"
)

type KcpReq struct {
	ICachedData;
	LEN uint16;
	UID uint32;
	RID uint32;
	BDY []byte;
}
func (me *KcpReq)ReadAt(conn net.Conn)(e error){
	defer func(){
		if err:=recover();err!=nil{
			e=errors.New(fmt.Sprint(err));
		}
	}();
	conn.SetReadDeadline(time.Now().Add(time.Second*5));
	if _,e:=io.ReadFull(conn,me.BDY[0:4]);e!=nil{
		return e;
	}
	if(binary.BigEndian.Uint32(me.BDY[0:4])!=12345){
		return errors.New("packet read error!");
	}
	if _,e:=io.ReadFull(conn,me.BDY[4:6]);e!=nil{
		return e;
	}
	me.LEN=binary.BigEndian.Uint16(me.BDY[4:6]);
	if _,e:=io.ReadFull(conn,me.BDY[6:me.LEN+4]);e!=nil{
		return e;
	}
	me.UID=binary.BigEndian.Uint32(me.BDY[6:10]);
	me.RID=binary.BigEndian.Uint32(me.BDY[10:14]);
	return nil;
}
func (me *KcpReq)GetALL()[]byte{
	return me.BDY;
}
func (me *KcpReq)GetLEN()uint16{
	return me.LEN;
}
func (me *KcpReq)GetUID()uint32{
	return me.UID;
}
func (me *KcpReq)GetRID()uint32{
	return me.RID;
}
func (me *KcpReq)GetRecvData()[]byte{
	return me.BDY[14:me.LEN+4];
}

type UdpReq struct {
	ICachedData;
	ADR net.Addr;
	LEN uint16;
	UID uint32;
	RID uint32;
	BDY []byte;
}
func (me *UdpReq)ReadAt(conn net.PacketConn)(err error){
	defer func(){
		if e:=recover();err!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}();
	_,a,e:=conn.ReadFrom(me.BDY);
	if e!=nil{
		return e;
	}
	me.LEN=binary.BigEndian.Uint16(me.BDY[4:6]);
	me.UID=binary.BigEndian.Uint32(me.BDY[6:10]);
	me.RID=binary.BigEndian.Uint32(me.BDY[10:14]);
	me.ADR=a;
	return nil;
}
func (me *UdpReq)GetALL()[]byte{
	return me.BDY;
}
func (me *UdpReq)GetLEN()uint16{
	return me.LEN;
}
func (me *UdpReq)GetUID()uint32{
	return me.UID;
}
func (me *UdpReq)GetRID()uint32{
	return me.RID;
}
func (me *UdpReq)GetRecvData()[]byte{
	return me.BDY[14:me.LEN+4];
}
func (me *UdpReq)GetAdr()net.Addr{
	return me.ADR;
}

