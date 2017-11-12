package room

import (
	"time"
	"github.com/sirupsen/logrus"
	"errors"
	"fmt"
	"runtime/debug"
)

func (me *Room1v1)room_kcp_recv_proc(plr *room_player){
	session:=plr.kcp_session;
	t:=time.Tick(time.Second*5);
	recv_flag:=false;
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
				logrus.Error(e);
				logrus.Error(fmt.Sprintf("%s",debug.Stack()));
			}
		}()
		for{
			select {
			case <-session.ChDie:
				return errors.New(fmt.Sprint("room ",me.rid,"session close_sig has called"));
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
				return errors.New(fmt.Sprint("room ",me.rid,"room close_sig has called"));
			}
		}

	}();
	if e!=nil{
		logrus.Error(e);
	}
	me.event_sig<-&kcp_session_closed{plr,session};
	session.Close(false);
}
