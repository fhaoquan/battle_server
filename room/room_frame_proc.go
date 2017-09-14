package room

import (
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"errors"
)

func (me *Room1v1)frame_proc(duration time.Duration){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	<-time.After(duration);
	me.event_sig<-&start_event{};
	frame:=0;
	t:=time.NewTicker(time.Millisecond*1000);
	f:=func()(run bool,err error){
		run=true;
		err=nil;
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		select {
		case <-me.close_sig:
			return false,nil;
		case <-t.C:
			me.event_sig<- &frame_event{uint32(frame)};
			frame++;
		}
		return ;
	};
	for{
		still_run,e:=f();
		if e!=nil{
			logrus.Error(e);
		}
		if !still_run{
			t.Stop();
			return;
		}
	}

}