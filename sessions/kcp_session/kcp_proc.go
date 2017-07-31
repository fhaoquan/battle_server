package kcp_session

import (
	"../../world"
	"../../utils"
	"net"
	"io"
	"encoding/binary"
	"errors"
)
func uuu(conn net.Conn,w *world.World){
	pool:=utils.NewMemoryPool(16,func()interface{}{
		return &kcp_packet{
			0,0,0,make([]byte,utils.MaxPktSize),
		}
	});
	recv_tkn:=make(chan interface{},1);
	main_tkn:=make(chan utils.IDataOwner,8);
	chk:=make([]byte,4);
	f1:=(func())(nil);
	f1=func(){
		dat:=pool.PopOne();
		if e:=func(tkn chan interface{})(error){
			defer func(){
				recv_tkn<-1
			}();
			<-recv_tkn;
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
			main_tkn<-dat;
			return nil;

		}(recv_tkn);e!=nil{
			return ;
		}
		go f1();



		<-main_tkn;
		buf:=dat.GetUserData().(*kcp_packet);
		r:=w.FindRoom(buf.r);
		r.OnKcp(buf.l,buf.u,buf.r,buf.b[10:buf.l])
		main_tkn<-1;

	};
	go f1();
}
