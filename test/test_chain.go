package test

import "../room"
type s_player_info struct {
	id int;
	name string;
}
func (me *s_player_info)GetPlayerID()uint32{
	return uint32(me.id);
}
func (me *s_player_info)GetPlayerName()string{
	return me.name;
}
type new_room_request_json struct{
	room_type int;
	room_players []s_player_info;
}
func TestNewRoom(){
	room.BuildRoom1v1(&s_player_info{1,"e"},&s_player_info{2,"e"});
}
