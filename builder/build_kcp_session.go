package builder

import (
	"../sessions/kcp_session"
	"../room"
	"../world"
	"net"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func BuildKcpSession(conn net.Conn,world *world.World){
	s:=kcp_session.NewSession(conn);
	f:=func(r *room.Room) {
		new(kcp_session.SendUtilErrorContext).
			WithSession(s).
			WithMsgPuller(
			func([]byte)int{
				return 0;
			}).
			WithErrHandle(
			func(err error){
				logrus.Error(err);
			}).
			SendUtilError();
	}
	go func(r *room.Room){
		new(kcp_session.RecvUtilErrorContext).
			WithSession(s).
			WithMsgPusher(
			func(data_user func(func(uid uint32,rid uint32,bdy []byte)bool))(error){
				if(r==nil){
					data_user(func(uid uint32, rid uint32, bdy []byte)bool{
						if r=world.FindRoom(uid);r!=nil{
							go f(r);
							return false;
						}
						return true;
					})
				}
				if(r!=nil){
					r.OnKcp(data_user);
					return nil;
				}
				return errors.New("cant find room");
			}).
			WithErrHandle(
			func(err error){
				logrus.Error(err);
			}).
			RecvUtilError();
	}(nil);
}