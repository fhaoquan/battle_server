package command

import (
	"../battle"
	"fmt"
	"errors"
)

func (cmd *CommandContext)BroadcastBattleMovementData()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=cmd.udp_res_pool.Pop().(*udp_response);
	wtr:=&packet_encoder{
		res.bdy,
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	count:=0;
	cmd.base_room.GetBattle().ForEachUnitDo(func(u *battle.Unit)bool{
		wtr.write_unit_id(u.ID).
			write_unit_x(u.X).
			write_unit_y(u.Y).
			write_unit_speed(u.Speed).
			write_unit_face(u.Direction).
			write_unit_aiming_face(u.AimingFace);
		count++;
		return true;
	})
	ph1(uint8(count));
	res.len=uint16(wtr.pos);
	ph0(res.len-2);
	return res;
}
