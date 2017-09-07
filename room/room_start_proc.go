package room

import (
	"fmt"
	"time"
	"errors"
)

func (me *Room1v1)start_proc(){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	if e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		for {
			select {
			case _,ok:=<-me.close_sig:
				if !ok {
					return errors.New("start stoped by close signal");
				}
			case <-time.After(time.Second*5):
				return errors.New(fmt.Sprint("room ",me.rid, " wait player login timeout"));
			case event:=<-me.event_sig:
				me.on_event(event);
				if me.p1.kcp_session!=nil && me.p2.kcp_session!=nil{
					return nil;
				}
			}
		}
	}(); e!=nil{
		me.Close(e);
		return;
	}
	go me.udp_recv_proc(me.p1.udp_session.udp_conn.GetConn());
	go me.logic_proc();
	go me.frame_proc(time.Second*5);
}
