package room

import "github.com/pkg/errors"

type S_room_builder struct {
	r *Room;
}
func (me* S_room_builder)ID(v uint32)(*S_room_builder){
	me.r.id=v;
	return me;
}
func (me* S_room_builder)RouteCommand(cmd_id int,f func([]byte,*Room))(*S_room_builder){
	switch {
	case (0<=cmd_id&&cmd_id<len(me.r.cmd_handlers)):
		me.r.cmd_handlers[cmd_id]=f;
		return me;
	default:
		panic(errors.New("cmd_id error"));
	}
}
func (me* S_room_builder)RouteTimer(cmd_id int,f func(*Room))(*S_room_builder){
	return me;
}
func (me* S_room_builder)Players(i_all_players []interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
}){
	for i:=range i_all_players{
		me.r.players[i_all_players[i].GetPlayerID()]=&player{
			0,
			i_all_players[i].GetPlayerID(),
			i_all_players[i].GetPlayerName(),
			nil,
		}
	}
}
func (me* S_room_builder)Player(i_player_getter ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*S_room_builder){
	me.Players(i_player_getter);
	return me;
}
func (me* S_room_builder)Battle(battle i_battle)(*S_room_builder){
	me.r.battle=battle;
	return me;
}
func (me* S_room_builder)Build()(*Room){
	return me.r;
}
