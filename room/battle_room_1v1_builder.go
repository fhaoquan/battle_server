package room

import (
	"../udp_service"
	"../battle"
	"errors"
	"../command"
)

type BattleRoomBuilder struct {
	r *BattleRoom1v1;
}

func NewBattleRoomBuilder(r *BattleRoom1v1)(*BattleRoomBuilder){
	return &BattleRoomBuilder{r};
}
func (me* BattleRoomBuilder)RouteCommand(cmd_id int,f func([]byte)(interface{}))(*BattleRoomBuilder){
	switch {
	case (0<=cmd_id&&cmd_id<len(me.r.cmd_handlers)):
		me.r.cmd_handlers[cmd_id]=f;
		return me;
	default:
		panic(errors.New("cmd_id error"));
	}
}
func (me* BattleRoomBuilder)RouteTimer(cmd_id int,f func()(interface{}))(*BattleRoomBuilder){
	switch {
	case (0<=cmd_id&&cmd_id<len(me.r.cmd_handlers)):
		me.r.timer_handlers[cmd_id]=f;
		return me;
	default:
		panic(errors.New("cmd_id error"));
	}
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
func (me* BattleRoomBuilder)WithUDPSession(connection *udp_service.UdpConnection)(*BattleRoomBuilder){
	me.r.rid=(uint32)(connection.Addr.Port);
	me.r.p1.udp_session=&battle_udp_session{me.r.p1.uid,connection,nil};
	me.r.p2.udp_session=&battle_udp_session{me.r.p2.uid,connection,nil};
	return me;
}
func BuildRoom1v1(plrs ...interface{
	GetPlayerID()uint32;
	GetPlayerName()string;
})(*BattleRoom1v1,error){
	u,e:=udp_service.TheUDPConnManager.Pop();
	if(e!=nil){
		return nil,e;
	}
	c:=&command.NewCommandContext();
	r:=NewBattleRoomBuilder(&BattleRoom1v1{
		new_battle_room_base(battle.NewBattle()),
		nil,
		nil,
	}).
		WithPlayers(plrs[0],plrs[1]).
		WithUDPSession(u).
		RouteCommand(002,c.CreateUnit).
		RouteCommand(003,c.UnitAttackStart).
		RouteCommand(004,c.UnitAttackDone).
		RouteCommand(005,c.UpdateUnitMovement).
		RouteTimer(001,c.BroadcastBattleMovementData).r;
	c.SetRoom(r);
	return r,nil;
}