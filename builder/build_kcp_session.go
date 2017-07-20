package builder

import (
	"../sessions"
	"../room"
	"../world"
	"net"
	"time"
)
func BuildKcpSession(conn net.Conn,world *world.World){
	s:=sessions.NewSession(conn);
	c:=make(chan *room.Room,1);
	go func(){
		r:=room.NilRoom();
		sessions.NewReadLoop().
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
			sessions.NewSendLoop().
				WithSession(s).
				WithMsgGetter(func()[]byte{
					return nil;
				}).
				Do();
		case <-time.After(time.Minute):
			return;
		}
	}();
}