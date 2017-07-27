package udp_session

import (
	"net"
	"fmt"
	"../../utils"
	"sync"
	"io"
	"encoding/binary"
	"errors"
)

type Session struct {
	con net.PacketConn;
	buf []byte;
}
func (s *Session)ReadPacket(buf []byte)(net.Addr,error){
	l,a,e:=s.con.ReadFrom(s.buf);
	if e!=nil{
		return a,e;
	}
	if(binary.BigEndian.Uint32(s.buf[0:4])!=12345){
		return a,errors.New("packet read error!");
	}
	copy(buf,s.buf[4:l]);
	return a,nil;
}
func NewSession(port int)(*Session,error){
	if adr,err:=net.ResolveUDPAddr("udp",fmt.Sprint(":",port));err!=nil{
		return nil,err;
	}else if con,err:=net.ListenUDP("udp", adr);err!=nil{
		return nil,err;
	}else{
		con.SetWriteBuffer(utils.MaxPktSize*16);
		con.SetReadBuffer(utils.MaxPktSize*16);
		return &Session{
			con:con,
		},nil;
	}
}
