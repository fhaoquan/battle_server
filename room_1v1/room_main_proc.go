package room_1v1

import (
	"time"
	"errors"
	"../utils"
	"../gateway"
)

func (me *room)wait_players()(error){
	p1_started:=false;
	p2_started:=false;
	for{
		select {
		case <-me.close_sig:
			return errors.New("close signal checked");
		case <-time.After(time.Minute):
			return errors.New("wait player timeout");
		case event:=<-me.event_sig:
			switch event.(type){
			case *event_session_connected:
				me.on_event_session_connected(event.(*event_session_connected));
			case *event_session_closed:
				me.on_event_session_closeed(event.(*event_session_closed))
			}
		case pkt:=<-me.packet_chan:
			switch pkt.GetMsgBody()[0]{
			case utils.CMD_battle_start:
				if pkt.GetUID()==me.p1.uid{
					p1_started=true;
				}
				if pkt.GetUID()==me.p2.uid{
					p2_started=true;
				}
				if p1_started && p2_started{
					pkt.Return();
					return nil;
				}
			}
			pkt.Return();

		}
	}
}
func (me *room)loop_until_close(){
	t1:=time.NewTicker(time.Millisecond*50);
	t2:=time.NewTicker(time.Second);
	defer t1.Stop();
	defer t2.Stop();
	for{
		select {
		case <-me.close_sig:
			me.room_log_inf("close signal checked");
			return;
		case event:=<-me.event_sig:
			me.on_event(event);
		case pkt:=<-me.packet_chan:
			me.on_packet(pkt.GetUID(),pkt.GetMsgBody());
			pkt.Return();
		case <-t1.C:
			me.on_timer(time.Millisecond*50);
		case <-t2.C:
			me.on_timer(time.Second);
		}
	}
}
func (me *room)main_proc(){
	defer func() {
		gateway.DelReceiver(me.rid,me.p1.uid);
		gateway.DelReceiver(me.rid,me.p2.uid);
		me.room_log_inf("main proc exited");
		if e:=recover();e!=nil{
			me.room_log_err(e);
		}
	}()
	gateway.AddReceiver(me.rid,me.p1.uid,me);
	gateway.AddReceiver(me.rid,me.p2.uid,me);
	if e:=me.wait_players();e!=nil{
		me.room_log_err(e);
		return ;
	}
	me.handle_battle_result(0,me.the_battle.BroadcastBattleStart(0))
	me.handle_battle_result(0,me.the_battle.BroadcastBattleAll(0))
	me.loop_until_close();
}
