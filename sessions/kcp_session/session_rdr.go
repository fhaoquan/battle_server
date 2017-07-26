package kcp_session

import (
	"io"
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"fmt"
)
type onMsgFunc func(uid uint32,rid uint32,bdy []byte)
type read_loop_context struct {
	s *session;
	on_msg onMsgFunc;
}
func (me *read_loop_context)Do(){
	buf:=make([]byte,1024);
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		for{
			if e:=me.s.ReadPacket(buf);e==nil{
				l:=binary.BigEndian.Uint16(buf[0:2]);
				u:=binary.BigEndian.Uint32(buf[2:6]);
				r:=binary.BigEndian.Uint32(buf[6:10]);
				me.on_msg(u,r,buf[10:l+2]);
			}else{
				return e;
			}
		}
		return nil;
	}();
	me.s.con.Close();
	logrus.Error(err);
}
type i_WithMsgReceiverRtn interface{Do()}
func (me *read_loop_context)WithMsgReceiver(on_msg onMsgFunc)(i_WithMsgReceiverRtn){
	me.on_msg=on_msg;
	return me;
}
type i_WithSessionRtn interface{WithMsgReceiver(on_msg onMsgFunc)(i_WithMsgReceiverRtn)}
func (me *read_loop_context)WithSession(s *session)i_WithSessionRtn{
	return me;
}
type i_RdlpRtn interface {WithSession(*session)(i_WithSessionRtn)}
func NewReadLoop()i_RdlpRtn  {
	return &read_loop_context{}
}