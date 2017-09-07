package room

import (
	"../server/udp_server"
	"net"
	"fmt"
	"errors"
)

type udp_session struct {
	uid uint32;
	udp_conn *udp_server.UdpConnection;
	udp_addr net.Addr;
}
func (me *udp_session)Close(){
	if c:=me.udp_conn.GetConn();c!=nil{
		c.Close();
	}

	me.udp_conn.Return();
}
func (me *udp_session)Send(b []byte)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	if c:=me.udp_conn.GetConn();c!=nil{
		_,err=c.WriteTo(b,me.udp_addr);
		return ;
	}else{
		return errors.New("udp_conn is null");
	}
}