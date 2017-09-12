package room

import (
	"../battle"
	"errors"
)

type BattleRoomBuilder struct {
	r *Room1v1;
}

func NewBattleRoomBuilder(r *Room1v1)(*BattleRoomBuilder){
	return &BattleRoomBuilder{r};
}
func (me* BattleRoomBuilder)WithPlayers(i_player_getter ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*BattleRoomBuilder){
	if len(i_player_getter)!=2{
		panic(errors.New("1v1 room need 2 players"));
		return me;
	}
	me.r.p1=&room_player{
		i_player_getter[0].GetPlayerID(),
		i_player_getter[0].GetPlayerName(),
		nil,
		nil,
	}
	me.r.p2=&room_player{
		i_player_getter[1].GetPlayerID(),
		i_player_getter[1].GetPlayerName(),
		nil,
		nil,
	}
	return me;
}
func BuildRoom1v1(plrs ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*Room1v1,error){
	r:=NewBattleRoomBuilder(&Room1v1{
		new_base_room(battle.NewBattle()),
		nil,
		nil,
	}).WithPlayers(plrs[0],plrs[1]).r;
	return r,nil;
}