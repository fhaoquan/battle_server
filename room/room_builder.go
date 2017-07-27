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
func (me* S_room_builder)WithPlayers(i_player_getter ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*S_room_builder){
	for i:=range i_player_getter{
		me.r.players[i_player_getter[i].GetPlayerID()]=&player{
			0,
			i_player_getter[i].GetPlayerID(),
			i_player_getter[i].GetPlayerName(),
		}
	}
	return me;
}
func (me* S_room_builder)Build()(*Room){
	return me.r;
}
