package room

import (
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"errors"
	"runtime/debug"
)

func (me *Room1v1)frame_proc(duration time.Duration){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	<-time.After(duration);
	me.event_sig<-&start_event{1};
	frame:=0;
	t:=time.NewTicker(time.Millisecond*50);
	f:=func()(run bool,err error){
		run=true;
		err=nil;
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
				logrus.Error(e);
				logrus.Error(fmt.Sprintf("%s",debug.Stack()));
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