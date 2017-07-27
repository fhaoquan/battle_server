package udp_session

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net"
)

type SendUtilErrorContext struct {
	s *Session;
	msg_puller func([]byte)(net.Addr,int)
	err_handle func(err error);
}
func (context *SendUtilErrorContext)WithSession(s*Session)*SendUtilErrorContext{
	context.s=s;
	return context;
}
func (context *SendUtilErrorContext)WithMsgPuller(f func([]byte)(net.Addr,int))*SendUtilErrorContext{
	context.msg_puller=f;
	return context;
}
func (context *SendUtilErrorContext)WithErrHandle(f func(err error))*SendUtilErrorContext{
	context.err_handle=f;
	return context;
}
func (context *SendUtilErrorContext)SendUtilError(){
	buf:=make([]byte,1024);
	err:= func()(e error) {
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		for {
			if a,l:=context.msg_puller(buf);l>0{
				if _,e:=context.s.con.WriteTo(buf[:l],a);e!=nil{
					return e;
				}
			}
		}
		return nil;
	}();
	logrus.Error(err);
}
