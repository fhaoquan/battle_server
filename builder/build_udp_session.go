package builder

import (
	"../room"
	"../sessions/udp_session"
	"net"
	"github.com/sirupsen/logrus"
	"errors"
	"time"
)


func BuildUdpSession(start_port int)(func(r *room.Room)(uint32),error){
	f:=func(port int,conn net.PacketConn)func(r *room.Room)(uint32){
		return func(r *room.Room)(uint32){
			go func() {
				new(udp_session.SendUtilErrorContext).
					WithSession(s).
					WithMsgPuller(
					func([]byte)(net.Addr,int){
						return nil,0;
					}).
					WithErrHandle(
					func(err error){
						logrus.Error(err);
					}).
					SendUtilError();
			}();
			go func(){
				logrus.Fatal(udp_session.RecvUntilFalse(
					conn,
					func(addr net.Addr,len uint16,uid uint32,rid uint32,bdy []byte)bool {
						return r.OnUdp(addr,len,uid,rid,bdy)==nil;
					},
				))
				conn.Close();
			}();
			return uint32(port);
		}
	}
	t:=time.Now()
	for{
		if(time.Now().Sub(t).Minutes()>1){
			return nil,errors.New("cant listen udp");
		}
		if s,e:=udp_session.TryListen(start_port);e==nil{
			return f(start_port,s),nil;
		}
		start_port++;
		time.Sleep(time.Millisecond*100);
	}
	return nil,errors.New("cant listen udp");
}
