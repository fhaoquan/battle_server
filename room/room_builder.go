package room

import "github.com/pkg/errors"

type RoomBuilder struct {
	r *Room;
}
func (me* RoomBuilder)ID(v uint32)(*RoomBuilder){
	me.r.id=v;
	return me;
}
func (me* RoomBuilder)RouteCommand(cmd_id int,f func([]byte)(interface{}))(*RoomBuilder){
	switch {
	case (0<=cmd_id&&cmd_id<len(me.r.cmd_handlers)):
		me.r.cmd_handlers[cmd_id]=f;
		return me;
	default:
		panic(errors.New("cmd_id error"));
	}
}
func (me* RoomBuilder)RouteTimer(cmd_id int,f func()(interface{}))(*RoomBuilder){
	switch {
	case (0<=cmd_id&&cmd_id<len(me.r.cmd_handlers)):
		me.r.timer_handlers[cmd_id]=f;
		return me;
	default:
		panic(errors.New("cmd_id error"));
	}
}
func (me* RoomBuilder)WithPlayers(i_player_getter ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*RoomBuilder){
	for i:=range i_player_getter{
		me.r.players[i_player_getter[i].GetPlayerID()]=&Player{
			0,
			i_player_getter[i].GetPlayerID(),
			i_player_getter[i].GetPlayerName(),
		}
	}
	return me;
}
func (me* RoomBuilder)Build()(*Room){
	return me.r;
}
