package builder

import (
	"../sessions/kcp_session"
	"../room"
	"../world"
	"net"
	"time"
)
func BuildKcpSession(conn net.Conn,world *world.World){
	s:=kcp_session.NewSession(conn);
	c:=make(chan *room.Room,1);
	go func(){
		r:=room.NilRoom();
		kcp_session.NewReadLoop().
			WithSession(s).
			WithMsgReceiver(func(uid uint32,rid uint32,bdy []byte){
				if(r==nil){
					if r=world.FindRoom(uid);r!=nil{
						c<-r;
					}
				}
				if(r!=nil){
					r.OnPkt(uid,rid,bdy);
				}
			}).
			Do();
	}();
	go func(){
		select {
		case r,ok:=<-c:
			if(!ok||r==nil){
				return;
			}
			kcp_session.NewSendLoop().
				WithSession(s).
				WithMsgGetter(func(buf []byte)int{
					return 0;
				}).
				Do();
		case <-time.After(time.Minute):
			return;
		}
	}();
}