package room_1v1

import (
	"time"
	"sync"
	"../utils"
	"../battle"
	"../gateway"
	"../result_cache"
)

type room struct {
	gid         	string
	rid         	uint32
	start_timer 	time.Time
	lifecycle		time.Duration;
	win_score		int;
	the_battle  	*battle.Battle;
	event_sig   	chan interface{}
	close_sig   	chan interface{}
	packet_chan    	chan utils.I_REQ
	once_start  	*sync.Once
	once_close  	*sync.Once
	wait        	*sync.WaitGroup
	sudden_death 	time.Duration
	schedule_status	uint16;
	p1				*player;
	p2				*player;
}
func (me *room) Start(){
	me.once_start.Do(func() {
		go me.main_proc();
	})
}
func (me *room) NewSession(s gateway.Session){
	me.event_sig<-&event_session_connected{s};
}
func (me *room) DelSession(s gateway.Session){
	me.event_sig<-&event_session_closed{s};
}
func (me *room) GetGuid()string {
	return me.gid
}
func (me *room) SetID(v uint32) {
	me.rid = v
}
func (me *room) GetID() uint32 {
	return me.rid
}
func (me *room) GetBattle() *battle.Battle {
	return me.the_battle
}
func (me *room) Close(why error) {
	go me.once_close.Do(func() {
		me.room_log_inf(" will closed for :", why);
		close(me.close_sig)
		me.wait.Wait()
		res := new_room_result(me.gid)
		res.Result = append(res.Result, new_player_result(int(me.p1.uid), me.the_battle.ComputeResultScore(me.p1.uid)))
		res.Result = append(res.Result, new_player_result(int(me.p2.uid), me.the_battle.ComputeResultScore(me.p2.uid)))
		result_cache.CacheResult(res.Guid, res);
		del_room(me.rid);
		me.room_log_inf("room closed;key=",me.gid);
	})
}
