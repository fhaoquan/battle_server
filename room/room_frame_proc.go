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
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		for {
			select {
			case _,ok:=<-me.close_sig:
				if !ok {
					return nil;
				}
			default:
				time.Sleep(time.Millisecond*50);
				me.event_sig<- frame_event{uint32(frame)};
				frame++;
			}
		}
	}();
	if e!=nil{
		logrus.Error(e);
	}
}