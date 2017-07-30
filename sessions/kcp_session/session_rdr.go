package kcp_session

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"fmt"
	"../../utils"
	"net"
	"io"
)
func RecvUntilError(conn net.Conn,receiver func(uint16,uint32,uint32,[]byte)(error))(e error){
	defer func(){
		if err:=recover();err!=nil{
			e=errors.New(fmt.Sprint(err));
		}
	}();
	new:=true;
	chk:=make([]byte,4);
	buf:=make([]byte,utils.MaxPktSize)
	for{
		if _,e:=io.ReadFull(conn,chk);e!=nil{
			return e;
		}
		if(binary.BigEndian.Uint32(chk)!=123454321){
			return errors.New("packet read error!");
		}
		if _,e:=io.ReadFull(conn,buf[0:2]);e!=nil{
			return e;
		}
		l:=binary.BigEndian.Uint16(buf[0:2]);
		if _,e:=io.ReadFull(conn,buf[2:l]);e!=nil{
			return e;
		}
		u:=binary.BigEndian.Uint32(buf[2:6]);
		r:=binary.BigEndian.Uint32(buf[6:10]);
		if(new){

		}
		if e:=receiver(l,u,r,buf[10:l]);e!=nil{
			return errors.New(fmt.Sprintf("stoped for receiver error=",e))
		}
	}
}
func RunUntilRoomEnd()func(uint16,uint32,uint32,[]byte)(error){
	msg_c:=make(chan int,16);
	for{
		select {
		case <-msg_c:
		}
	}
	return func(uint16,uint32,uint32,[]byte)(error){
		return nil;
	}
}
func SendUtilAnyError()(error){

}