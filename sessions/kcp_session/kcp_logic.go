package kcp_session

import (
	"../../utils"
	"../../room"
	"../../world"
	"net"
	"io"
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
)
func startSendProc(room *room.Room){

}
func startMainProc(room *room.Room,msg chan utils.IDataOwner,sig chan interface{}){
	err:= func()(error){
		for{
			select {
			case data:=<-msg:
				pkt:=data.GetUserData().(*kcp_packet)
				room.OnKcp(pkt.l,pkt.u,pkt.r,pkt.b[10:pkt.l]);
				data.Return();
			case <-sig:
				return errors.New("stoped!");
			}
		}
	}()
	logrus.Error(err);
}
func StartRecvProc(conn net.Conn,w *world.World){
	pool:=utils.NewMemoryPool(16,func()interface{}{
		return &kcp_packet{
			0,0,0,make([]byte,utils.MaxPktSize),
		}
	});
	room:=(*room.Room)(nil);
	chk:=make([]byte,4);
	msg:=make(chan utils.IDataOwner,16);
	sig:=make(chan interface{},1);

	go func(f func(dat utils.IDataOwner)(error)){
		err:=error(nil);
		for{
			dat:=pool.PopOne();
			if err=f(dat);err!=nil{
				dat.Return();
				break;
			}
		}
		close(msg);
		sig<-1;
		logrus.Error(err);
	}(
		func(dat utils.IDataOwner)(error){
			buf:=dat.GetUserData().(*kcp_packet);
			if _,e:=io.ReadFull(conn,chk);e!=nil{
				return e;
			}
			if(binary.BigEndian.Uint32(chk)!=123454321){
				return errors.New("packet read error!");
			}
			if _,e:=io.ReadFull(conn,buf.b[0:2]);e!=nil{
				return e;
			}
			buf.l=binary.BigEndian.Uint16(buf.b[0:2]);
			if _,e:=io.ReadFull(conn,buf.b[2:buf.l]);e!=nil{
				return e;
			}
			buf.u=binary.BigEndian.Uint32(buf.b[2:6]);
			buf.r=binary.BigEndian.Uint32(buf.b[6:10]);
			if(room==nil){
				room=w.FindRoom(buf.u);
				if(room!=nil){
					startSendProc(room);
					startMainProc(room,msg,sig);
				}else{
					return errors.New("packet read error!");
				}
			}
			if(room!=nil){
				msg<-dat;
			}
			return nil;
		},
	);

}
