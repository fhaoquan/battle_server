package room_1v1

import (
	"time"
	"errors"
	)

func (me *room)on_timer_frame(){
	remaining_time := me.lifecycle-time.Now().Sub(me.start_timer);
	me.handle_battle_result(0, me.the_battle.BroadcastBattleMovementData(
		me.p1.uid,
		me.p2.uid,
		me.schedule_status,
		uint16(remaining_time.Seconds())));
	me.handle_battle_result(0, me.the_battle.BroadcastBattleMovementData(
		me.p2.uid,
		me.p1.uid,
		me.schedule_status,
		uint16(remaining_time.Seconds())));
}
func (me *room)on_timer_second(){
	s1 := me.the_battle.ComputeResultScore(me.p1.uid)
	s2 := me.the_battle.ComputeResultScore(me.p2.uid)
	switch time_span := time.Now().Sub(me.start_timer); {
	case time_span < me.sudden_death:
		me.schedule_status=1;
		switch {
		case s1 >= me.win_score:
			me.handle_battle_result(0, me.the_battle.BroadcastBattleEnd(me.p1.uid))
			me.Close(errors.New("the battle complated!"))
		case s2 >= me.win_score:
			me.handle_battle_result(0, me.the_battle.BroadcastBattleEnd(me.p2.uid))
			me.Close(errors.New("the battle complated!"))
		}
	case time_span < me.lifecycle:
		me.schedule_status=2;
		switch {
		case s1 > s2:
			me.handle_battle_result(0, me.the_battle.BroadcastBattleEnd(me.p1.uid))
			me.Close(errors.New("the battle complated!"))
		case s1 < s2:
			me.handle_battle_result(0, me.the_battle.BroadcastBattleEnd(me.p2.uid))
			me.Close(errors.New("the battle complated!"))
		}
	default:
		me.schedule_status=3;
		me.handle_battle_result(0, me.the_battle.BroadcastBattleEnd(0))
		me.Close(errors.New("the battle complated!"))
		return
	}
}

func (me *room)on_timer(timer_type time.Duration){
	switch timer_type{
	case time.Millisecond*50:
		me.on_timer_frame();
	case time.Second:
		me.on_timer_second();
	}
}
