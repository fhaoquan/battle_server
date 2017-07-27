package proc

import (
	"../utils"
	"../world"
	"../room"
	"net"
	"fmt"
	"io"
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
)
func empty_room()*room.Room{
	return nil;
}
func recv_pkt(conn net.Conn,buf []byte)(error){
	if _,e:=io.ReadFull(conn,buf[0:4]);e!=nil{
		return e;
	}
	if(binary.BigEndian.Uint32(buf[0:4])!=12345){
		return errors.New("packet read error!");
	}
	if _,e:=io.ReadFull(conn,buf[4:6]);e!=nil{
		return e;
	}
	l:=binary.BigEndian.Uint16(buf[4:6]);
	if _,e:=io.ReadFull(conn,buf[6:l]);e!=nil{
		return e;
	}
	return nil;
}

func recv_proc(conn net.Conn,w *world.World){
	buf:=make([]byte,utils.MaxPktSize);
	tar:=empty_room();
	if err:=recv_pkt(conn,buf);err!=nil{
		logrus.Error(err);
		return ;
	}
	l:=binary.BigEndian.Uint16(buf[4:6]);
	u:=binary.BigEndian.Uint32(buf[6:10]);
	r:=binary.BigEndian.Uint32(buf[10:14]);
	if tar==nil{
		if tar=w.FindRoom(r);tar==nil{
			logrus.Error("can not find room with id=",r);
			return ;
		}
	}
	tar.OnPkt(u,r,buf[14:l]);
}
func main_proc(*room.Room){

}