package room

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"errors"
)

func (me *Room1v1)logic_proc(){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	f:=func()(run bool,err error){
		run=true;
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		select {
		case <-me.close_sig:
			return false,nil;
		case kcp_msg:=<-me.kcp_chan:
			me.on_kcp_message(kcp_msg);
		case udp_msg:=<-me.udp_chan:
			me.on_udp_message(udp_msg);
		case event:=<-me.event_sig:
			me.on_event(event);
		}
		return true,nil;
	}
	for{
		still_run,e:=f();
		if e!=nil{
			logrus.Error(e);
		}
		if !still_run{
			return ;
		}
	}

}