package player

import (
	"../utils"
	"net"
	"fmt"
	"io"
	"errors"
	"encoding/binary"
)
type KcpReq struct {
	utils.ICachedData;
	l uint16;
	u uint32;
	r uint32;
	b []byte;
}
func (me *KcpReq)ReadAt(conn net.Conn)(e error){
	defer func(){
		if err:=recover();err!=nil{
			e=errors.New(fmt.Sprint(err));
		}
	}();
	if _,e:=io.ReadFull(conn,me.b[0:4]);e!=nil{
		return e;
	}
	if(binary.BigEndian.Uint32(me.b[0:4])!=12345){
		return errors.New("packet read error!");
	}
	if _,e:=io.ReadFull(conn,me.b[4:6]);e!=nil{
		return e;
	}
	me.l=binary.BigEndian.Uint16(me.b[4:6]);
	if _,e:=io.ReadFull(conn,me.b[6:me.l+4]);e!=nil{
		return e;
	}
	me.u=binary.BigEndian.Uint32(me.b[6:10]);
	me.r=binary.BigEndian.Uint32(me.b[10:14]);
	return nil;
}
func (me *KcpReq)GetALL()[]byte{
	return me.b;
}
func (me *KcpReq)GetLEN()uint16{
	return me.l;
}
func (me *KcpReq)GetUID()uint32{
	return me.u;
}
func (me *KcpReq)GetRID()uint32{
	return me.r;
}
func (me *KcpReq)GetRecvData()[]byte{
	return me.b[14:me.l+4];
}