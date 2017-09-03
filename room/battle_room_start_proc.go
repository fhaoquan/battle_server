package room

import (
	"../udp_service"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"errors"
)

func (me *BattleRoom1v1)start_proc(){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	if u,e:=udp_service.TheUDPConnManager.Pop();e!=nil{
		logrus.Error(e);
		return;
	}else{
		me.p1.udp_session=&battle_udp_session{me.p1.uid,u,nil};
		me.p2.udp_session=&battle_udp_session{me.p2.uid,u,nil};
	}
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
			case time.After(time.Second*30):
				return errors.New("wait player login timeout");
			case event:=<-me.event_sig:
				me.on_event(event);
				if me.p1.kcp_session!=nil && me.p2.kcp_session!=nil{
					return nil;
				}
			}
		}
	}(); e!=nil{
		logrus.Error(e);
		return;
	}
	go me.udp_recv_proc(me.p1.udp_session.udp_conn);
	go me.logic_proc();
	go me.frame_proc();
}
