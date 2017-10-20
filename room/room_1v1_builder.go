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
	GetUnits()[]battle.Unit;
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
	for i,_:=range i_player_getter[0].GetUnits(){
		me.r.the_battle.CreateUnitDo(func(unit *battle.Unit) {
			unit.SetAll(&i_player_getter[0].GetUnits()[i]);
		})
	}
	for i,_:=range i_player_getter[1].GetUnits(){
		me.r.the_battle.CreateUnitDo(func(unit *battle.Unit) {
			unit.SetAll(&i_player_getter[1].GetUnits()[i]);
		})
	}
	return me;
}
func BuildRoom1v1(plrs ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
	GetUnits()[]battle.Unit;
})(*Room1v1,error){
	r:=NewBattleRoomBuilder(&Room1v1{
		new_base_room(battle.NewBattle()),
		nil,
		nil,
	}).WithPlayers(plrs[0],plrs[1]).r;
	return r,nil;
}