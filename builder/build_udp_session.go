package builder

import (
	"../room"
	"../sessions/udp_session"
	"net"
	"github.com/sirupsen/logrus"
	"errors"
)


func BuildUdpSession(start_port int)(interface{}){
	f:=func(port int,s *udp_session.Session)func(r *room.Room)(uint32){
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
				new(udp_session.RecvUtilErrorContext).
					WithSession(s).
					WithMsgPusher(
					func(f func(func(addr net.Addr,uid uint32,rid uint32,bdy []byte)bool))(error){
						r.OnUdp(f);
						return nil;
					}).
					WithErrHandle(
					func(err error){
						logrus.Error(err);
					}).
					RecvUtilError();
			}();
			return uint32(port);
		}
	}
	for i:=0;i<1000;i++{
		if s,e:=udp_session.NewSession(start_port+i);e==nil{
			return f(start_port+i,s);
		}
	}
	return errors.New("cant listen udp");
}
