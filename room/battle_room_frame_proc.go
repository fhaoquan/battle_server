package room

import (
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"errors"
)

func (me *BattleRoom1v1)frame_proc(){
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
			case _,ok:=<-me.close_sig:
				if !ok {
					return nil;
				}
			default:
				time.Sleep(time.Millisecond*50);
				me.event_sig<-0;
			}
		}
	}();
	if e!=nil{
		logrus.Error(e);
	}
}