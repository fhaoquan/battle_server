package room

import (
	"time"
	"github.com/sirupsen/logrus"
	"errors"
	"fmt"
)

func (me *Room1v1)room_kcp_recv_proc(p *room_player){
	session:=p.kcp_session;
	t:=time.Tick(time.Second*5);
	recv_flag:=false;
	e:=func()(e error){
		for{
			select {
			case p,ok:=<-session.ChRecv:
				if ok{
					recv_flag=true;
					me.kcp_chan<-p;
				}else{
					return errors.New(fmt.Sprint("room ",me.rid," session.ChRecv closed"));
				}
			case <-t:
				if !recv_flag{
					return errors.New(fmt.Sprint("room ",me.rid," session timeout"));
				}else{
					recv_flag=false;
				}
			case <-me.close_sig:
				return errors.New(fmt.Sprint("room ",me.rid," close_sig has called"));
			}
		}

	}();
	if e!=nil{
		logrus.Error(e);
	}
	me.event_sig<-&kcp_session_closed{p,session};
	session.Close(false);
}
