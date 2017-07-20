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
	r:=room.NilRoom();
	c:=make(chan int,1);
	go sessions.NewReadLoop().
		WithSession(s).
		WithMsgReceiver(func(uid uint32,rid uint32,bdy []byte){
			if(r==nil){
				if r=world.FindRoom(uid);r!=nil{
					c<-1;
				}
			}
			if(r!=nil){
				r.OnPkt(uid,rid,bdy);
			}
		}).
		Do();
	go func(){
		select {
		case <-c:
			close(c);
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