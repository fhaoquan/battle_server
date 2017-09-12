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
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		for {
			select {
			case <-me.close_sig:
				return nil;
			case kcp_msg:=<-me.kcp_chan:
				me.on_kcp_message(kcp_msg);
			case udp_msg:=<-me.udp_chan:
				me.on_udp_message(udp_msg);
			case event:=<-me.event_sig:
				me.on_event(event);
			}
		}
	}();
	if e!=nil{
		logrus.Error(e);
	}
}