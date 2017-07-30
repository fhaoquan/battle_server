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
	f1:=func(r *room.Room,u uint32) {
		logrus.Fatal(kcp_session.SendUtilError(
			conn,
			func(bytes []byte)int {
				return -1;
			},
		))
		conn.Close();
	}
	f2:=func(r *room.Room){
		logrus.Fatal(kcp_session.RecvUntilError(
			conn,
			func(len uint16,uid uint32,rid uint32,bdy []byte)error{
				if(r==nil){
					if r=world.FindRoom(rid);r==nil{
						return errors.New(fmt.Sprintf("cant find room id=",uid));
					}else{
						go f1(r,uid);
					}
				}
				return r.OnKcp(len,uid,rid,bdy);
			},
		))
		conn.Close();
	}
	go f2(nil);
}