package kcp_session

import (
	"net"
	"io"
	"encoding/binary"
	"errors"
)

type session struct {
	head32 []byte;
	con net.Conn;
}

func (s *session)ReadPacket(buf []byte)(error){
	if _,e:=io.ReadFull(s.con,s.head32);e!=nil{
		return e;
	}
	if(binary.BigEndian.Uint32(s.head32)!=12345){
		return errors.New("packet read error!");
	}
	if _,e:=io.ReadFull(s.con,buf[0:2]);e!=nil{
		return e;
	}
	l:=binary.BigEndian.Uint16(buf[0:2]);
	if _,e:=io.ReadFull(s.con,buf[2:l]);e!=nil{
		return e;
	}
	return nil;
}

func NewSession(con net.Conn)(s *session){
	return &session{
		con:con,
	}
}
