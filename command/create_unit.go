package command

import (
	"../battle"
	"fmt"
	"errors"
)

func (cmd *CommandContext)CreateUnit(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=cmd.kcp_res_pool.Pop().(*kcp_response);
	res.broadcast=true;
	wtr:=&packet_encoder{
		res.bdy,
		0,
	}
	rdr:=&packet_decoder{
		data:data,
		pos:0,
	}
	f:=wtr.get_uint16_placeholder();
	count:=rdr.read_unit_count();
	wtr.write_unit_count(count);
	for i:=0;i<(int)(count);i++{
		cmd.base_room.GetBattle().CreateUnitDo(r.read_unit_id(), func(u *battle.Unit) {
			u.Camps=rdr.read_unit_camps();
			u.HP=rdr.read_unit_hp();
			u.X=rdr.read_unit_location_x();
			u.Y=rdr.read_unit_location_y();
			u.Speed=rdr.read_unit_speed();
			u.Direction=rdr.read_unit_face();
			u.AimingFace=rdr.read_unit_aiming_face();
			wtr.write_unit_id(u.ID)
			wtr.write_uint8(u.Camps)
			wtr.write_uint16(u.HP);
			wtr.write_uint16(u.X);
			wtr.write_uint16(u.Y);
			wtr.write_uint16(u.Speed);
			wtr.write_uint16(u.Direction);
			wtr.write_uint16(u.AimingFace);
		})
	}
	res.len=uint16(wtr.pos);
	f(res.len);
	return res;
}