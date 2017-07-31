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
func startSendProc(rid uint32,uid uint32,w *world.World)func(conn net.Conn,sig chan interface{}){
	r:=w.FindRoom(rid);
	if(r==nil){
		return nil;
	}
	p:=r.GetPlayer(uid);
	if(p==nil){
		return nil;
	}
	return func(conn net.Conn,sig chan interface{}){
		msg:=make(chan []byte,16);
		p.SetKcpSender(msg);
		for{
			select {
			case pkt:=<-msg:
				conn.Write(pkt);
			case <-sig:
				return ;
			}
		}
		p.SetKcpSender(nil);
	}
}
func startMainProc(rid uint32,w *world.World)func(msg chan utils.IDataOwner,sig chan interface{}){
	r:=w.FindRoom(rid);
	if(r==nil){
		return nil;
	}
	return func(msg chan utils.IDataOwner,sig chan interface{}){
		go func(on_msg func(data utils.IDataOwner)){
			for{
				select {
				case data:=<-msg:
					on_msg(data);
				case <-sig:
					return;
				}
			}
			logrus.Error("kcp session stoped");
		}(
			func(data utils.IDataOwner){
				defer func(){
					recover();
					data.Return();
				}();
				pkt:=data.GetUserData().(*kcp_packet)
				r.OnKcp(pkt.l,pkt.u,pkt.r,pkt.b[10:pkt.l]);
			},
		);
	}
}
func StartRecvProc(conn net.Conn,w *world.World){
	pool:=utils.NewMemoryPool(16,func()interface{}{
		return &kcp_packet{
			0,0,0,make([]byte,utils.MaxPktSize),
		}
	});
	chk:=make([]byte,4);
	msg:=make(chan utils.IDataOwner,16);
	sig:=make(chan interface{},1);
	new:=true;

	go func(read func(dat utils.IDataOwner)(error)){
		err:=error(nil);
		for{
			dat:=pool.PopOne();
			if err=read(dat);err!=nil{
				dat.Return();
				break;
			}
		}
		close(msg);
		sig<-1;
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
			if(new){
				f1:=startMainProc(buf.r,w);
				if(f1==nil){
					return errors.New("can not start main proc");
				}
				f2:=startSendProc(buf.r,buf.u,w);
				if(f2==nil){
					return errors.New("can not start send proc");
				}
				go f1(msg,sig);
				go f2(conn,sig);
			}
			msg<-dat;
			return nil;
		},
	);

}
