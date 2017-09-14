package room

import (
	"fmt"
	"time"
	"errors"
	"../server/udp_server"
)

func (me *Room1v1)start_proc(){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	t:=time.NewTicker(time.Millisecond*1000);
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		for {
			select {
			case <-me.close_sig:
				return errors.New("start stoped by close signal");
			case <-time.After(time.Second*60):
				return errors.New(fmt.Sprint("room ",me.rid, " wait player login timeout"));
			case <-t.C:
				me.on_event(&start_event{0});
			case event:=<-me.event_sig:
				me.on_event(event);
				if me.p1.kcp_session!=nil && me.p2.kcp_session!=nil{
				//if me.p1.kcp_session!=nil && me.p2.kcp_session!=nil{
					return nil;
				}
			}
		}
	}();
	t.Stop();
	if e!=nil{
		me.Close(e);
		return;
	}
	udp_server.UdpSlot[me.rid-10000]=me.udp_chan;
	go me.logic_proc();
	go me.frame_proc(time.Second);
}
