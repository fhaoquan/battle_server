package kcp_session

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type message_getter func([]byte)int
type send_loop_context struct {
	s *session;
	pop message_getter;
}

func (me *send_loop_context)Do(){
	buf:=make([]byte,1024);
	err:= func()(e error) {
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		for {
			if l:=me.pop(buf);l>0{
				if _,e:=me.s.con.Write(buf[:l]);e!=nil{
					return e;
				}
			}
		}
		return nil;
	}();
	logrus.Error(err);
}
type i_WithMsgGetterRtn interface{Do()}
func (me *send_loop_context)WithMsgGetter(getter message_getter)(i_WithMsgGetterRtn){
	return me;
}

type i_send_loop_context_WithSessionRtn interface{WithMsgGetter(getter message_getter)(i_WithMsgGetterRtn)}
func (me *send_loop_context)WithSession(s *session)i_send_loop_context_WithSessionRtn{
	return me;
}
type i_SdlpRtn interface {WithSession(*session)(i_send_loop_context_WithSessionRtn)}
func NewSendLoop()i_SdlpRtn  {
	return &send_loop_context{}
}