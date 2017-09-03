package room

import (
	"../udp_service"
	"net"
	"fmt"
	"errors"
)

type battle_udp_session struct {
	uid uint32;
	udp_conn *udp_service.UdpConnection;
	udp_addr net.Addr;
}
func (me *battle_udp_session)Send(b []byte)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	_,err=me.udp_conn.WriteTo(b,me.udp_addr);
	return ;
}