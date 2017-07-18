package room

import (
	"../command"
	"../battle"
	"../utils"
)

type unit_data_setter battle.Unit

func(me *unit_data_setter)SetLocation(x uint16,y uint16){
	(*battle.Unit)(me).X=x;
	(*battle.Unit)(me).Y=y;
}
func(me *unit_data_setter)SetMovement(speed uint16,face uint16,aiming_face uint16){
	(*battle.Unit)(me).Speed=speed;
	(*battle.Unit)(me).Direction=face;
	(*battle.Unit)(me).AimingFace=aiming_face;
}
func(me *unit_data_setter)SetCamps(v uint8){

}
func(me *unit_data_setter)SetAttackPower(v uint16){

}
func(me *unit_data_setter)GetAttackPower()uint16{
	return me.AttackPower;
}
func(me *unit_data_setter)GetHP()uint16{
	return me.HP;
}
func(me *unit_data_setter)GetID()uint16{
	return me.ID;
}
func(me *unit_data_setter)GetX()uint16{
	return me.X;
}
func(me *unit_data_setter)GetY()uint16{
	return me.Y;
}
func(me *unit_data_setter)GetSpeed()uint16{
	return me.Speed;
}
func(me *unit_data_setter)GetFace()uint16{
	return me.Direction;
}
func(me *unit_data_setter)GetAimingFace()uint16{
	return me.AimingFace;
}
func(me *unit_data_setter)SetHP(v uint16){

}
func(me *unit_data_setter)SetStat(v uint16){

}

type message_send_helper Room


func build_command_002_handle_func(b *battle.Battle)(func(d []byte,r *Room)){
	return func(d []byte,r *Room){
		command.UpdateUnitMovement(d,func(id uint16)(command.IUnitMovementDataSetter){
			return (*unit_data_setter)(b.FindUnit(id));
		})
	}
}
func build_command_003_handle_func(b *battle.Battle)(func(d []byte,r *Room)){
	return func(d []byte,r *Room){
		command.CreateUnit(d, func(id uint16) command.IUnitCreateDataSetter {
			return (*unit_data_setter)(b.NewUnit(id));
		})
	}
}
func build_command_004_handle_func(b *battle.Battle)(func(d []byte,r *Room)){
	s_buf:=make([]byte,utils.MaxPktSize);
	return func(d []byte,r *Room){
		command.UnitAttackDone(
			d,
			func(id uint16)command.IUnitAttackHelper{
				return (*unit_data_setter)(b.FindUnit(id));
			},
			s_buf,
			func(data []byte,len int) {
				r.Broadcast(data,len);
			},
		)
	}
}

func build_command_005_handle_func(b *battle.Battle)(func(d []byte,r *Room)){
	return func(d []byte,r *Room){
		command.UnitAttackStart(d, func(data []byte,len int) {
			r.Broadcast(data,len);
		})
	}
}

func build_timer_001_handle_func(b *battle.Battle)(func(r *Room)){
	s_buf:=make([]byte,utils.MaxPktSize);
	return func(r *Room){
		command.BroadcastBattleData(
			s_buf,
			func(f func(getter command.IUnitMovementDataGetter)) {
				units:=b.AllUnit();
				for i:=range units{
					f((*unit_data_setter)(units[i]));
				}
				f(nil);
			},
			func(data []byte,len int) {
				r.Broadcast(data,len);
			},
		)
	};
}

func NewRoom1v1()(*S_room_builder){
	b:=battle.NewBattle();
	return NewRoom().
		Battle(b).
		RouteCommand(0,nil).
		RouteCommand(1,nil).
		RouteCommand(2,build_command_002_handle_func(b)).
		RouteCommand(3,build_command_003_handle_func(b)).
		RouteCommand(4,build_command_004_handle_func(b)).
		RouteCommand(5,build_command_005_handle_func(b)).
		RouteTimer(0,nil).
		RouteTimer(1,build_timer_001_handle_func(b));
}
