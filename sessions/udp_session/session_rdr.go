package udp_session

import (
	"../../utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"encoding/binary"
	"github.com/pkg/errors"
)

type onMsgFunc func(addr *net.Addr,uid uint32,rid uint32,bdy []byte)
type read_loop_context struct {
	s *session;
	on_msg onMsgFunc;
}

func (context *read_loop_context)Do(){
	buf:=make([]byte,utils.MaxPktSize);
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
			context.s.con.Close();
		}();
		for{
			len,addr,err:=context.s.con.ReadFrom(buf);
			if(err!=nil){
				return err;
			}
			if l:=binary.BigEndian.Uint16(buf[0:2]);(int(l+2)!=len){
				return errors.New("pkt received len err");
			}
			u:=binary.BigEndian.Uint32(buf[2:6]);
			r:=binary.BigEndian.Uint32(buf[6:10]);
			context.on_msg(&addr,u,r,buf[10:len]);
		}
	}();
	logrus.Error(err);
}
type iDo interface{Do()}
func (context *read_loop_context)WithReceiver(f onMsgFunc)(iDo){
	context.on_msg=f;
	return context;
}
type iWithReceiver interface{WithReceiver(f onMsgFunc)(iDo)}
func (context *read_loop_context)WithSession(s *session)(iWithReceiver){
	context.s=s;
	return context;
}
type iWithSession interface{WithSession(s *session)(iWithReceiver)}
func NewReadLoop()iWithSession{
	return &read_loop_context{}
}
