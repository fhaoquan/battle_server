package kcp_session

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"fmt"
	"../../utils"
)

type cached_kcp_empty_packet struct {
	l uint16;
	u uint32;
	r uint32;
	b []byte;
	o utils.ICachedData;
}
func (me *cached_kcp_empty_packet)Clear(){
	me.l=0;
	me.u=0;
	me.r=0;
	me.o=nil;
}
func (me *cached_kcp_empty_packet)use_data_util_true(f func(uid uint32,rid uint32,bdy []byte)bool){
	if me.o.IsReturned(){
		return;
	}
	if(f(me.u,me.r,me.b[10:me.l])){
		me.o.Return();
	};
}

type RecvUtilErrorContext struct {
	s *Session;
	msg_pusher func(func(func(uid uint32,rid uint32,bdy []byte)bool))(error)
	err_handle func(err error);
}
func (context *RecvUtilErrorContext)WithSession(s*Session)*RecvUtilErrorContext{
	context.s=s;
	return context;
}
func (context *RecvUtilErrorContext)WithMsgPusher(f func(func(func(uid uint32,rid uint32,bdy []byte)bool))(error))*RecvUtilErrorContext{
	context.msg_pusher=f;
	return context;
}
func (context *RecvUtilErrorContext)WithErrHandle(f func(err error))*RecvUtilErrorContext{
	context.err_handle=f;
	return context;
}
func (context *RecvUtilErrorContext)RecvUtilError(){
	mpl:=utils.NewMemoryPool(16,func(i utils.ICachedData)utils.IUserData{
		return &cached_kcp_empty_packet{
			0,0,0,make([]byte,utils.MaxPktSize),i,
		}
	})
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		for{
			buf:=mpl.Get().GetUserData().(*cached_kcp_empty_packet);
			if e:=context.s.ReadPacket(buf.b);e!=nil{
				return e;
			}
			buf.l=binary.BigEndian.Uint16(buf.b[0:2]);
			buf.u=binary.BigEndian.Uint32(buf.b[2:6]);
			buf.r=binary.BigEndian.Uint32(buf.b[6:10]);
			if e:=context.msg_pusher(buf.use_data_util_true);e!=nil{
				return e;
			}
		}
		return nil;
	}();
	context.s.con.Close();
	logrus.Error(err);
}