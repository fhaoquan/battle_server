package udp_session

import (
	"../../utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"encoding/binary"
	"github.com/pkg/errors"
)

type cached_udp_empty_packet struct {
	a net.Addr;
	l uint16;
	u uint32;
	r uint32;
	b []byte;
	o utils.ICachedData;
}
func (me *cached_udp_empty_packet)Clear(){
	me.a=nil;
	me.l=0;
	me.u=0;
	me.r=0;
	me.o=nil;
}
func (me *cached_udp_empty_packet)use_data_util_true(f func(adr net.Addr,uid uint32,rid uint32,bdy []byte)bool){
	if me.o.IsReturned(){
		return;
	}
	if(f(me.a,me.u,me.r,me.b[10:me.l])){
		me.o.Return();
	}
}

type RecvUtilErrorContext struct {
	s *Session;
	msg_pusher func(func(func(addr net.Addr,uid uint32,rid uint32,bdy []byte)bool))(error)
	err_handle func(err error);
}
func (context *RecvUtilErrorContext)WithSession(s*Session)*RecvUtilErrorContext{
	context.s=s;
	return context;
}
func (context *RecvUtilErrorContext)WithMsgPusher(f func(func(func(addr net.Addr,uid uint32,rid uint32,bdy []byte)bool))(error))*RecvUtilErrorContext{
	context.msg_pusher=f;
	return context;
}
func (context *RecvUtilErrorContext)WithErrHandle(f func(err error))*RecvUtilErrorContext{
	context.err_handle=f;
	return context;
}

func (context *RecvUtilErrorContext)RecvUtilError(){
	mpl:=utils.NewMemoryPool(16,func(i utils.ICachedData)utils.IUserData{
		return &cached_udp_empty_packet{
			nil,0,0,0,make([]byte,utils.MaxPktSize),i,
		}
	})
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
			context.s.con.Close();
		}();
		for{
			buf:=mpl.Get().GetUserData().(*cached_udp_empty_packet);
			addr,err:=context.s.ReadPacket(buf.b);
			if(err!=nil){
				logrus.Error(err);
				continue;
			}
			buf.a=addr;
			buf.l=binary.BigEndian.Uint16(buf.b[0:2]);
			buf.u=binary.BigEndian.Uint32(buf.b[2:6]);
			buf.r=binary.BigEndian.Uint32(buf.b[6:10]);
			if e:=context.msg_pusher(buf.use_data_util_true);e!=nil{
				return;
			}
		}
	}();
	logrus.Error(err);
}
