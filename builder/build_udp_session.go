package builder

import (
	"../room"
	"../sessions/udp_session"
	"net"
)

func BuildUdpSession(start_port int,room *room.Room)(int){
	for i:=0;i<1000;i++{
		if s,e:=udp_session.NewSession(start_port+i);e==nil{
			go func(){
				udp_session.NewReadLoop().
					WithSession(s).
					WithReceiver(func(addr *net.Addr,uid uint32,rid uint32,bdy []byte){
					room.OnPkt(uid,rid,bdy);
				}).Do();
			}()

		}
	}
	return start_port;
}
