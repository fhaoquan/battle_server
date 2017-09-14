package battle

import (
	"fmt"
	"errors"
	"../utils"
)

func (context *Battle)CreateUnit(who uint32,data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	rdr:=&packet_decoder{
		data:data,
		pos:0,
	}
	ph1:=wtr.get_uint16_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=rdr.read_unit_count();
	wtr.write_unit_count(count);
	for i:=0;i<(int)(count);i++{
		context.CreateUnitDo(rdr.read_unit_id(), func(u *Unit) {
			u.Owner=who;
			u.Type=rdr.read_uint8();
			u.Camps=rdr.read_unit_camps();
			u.Level=rdr.read_uint8();
			u.Ranks=rdr.read_uint8();
			u.HP=rdr.read_unit_hp();
			u.X=rdr.read_unit_location_x();
			u.Y=rdr.read_unit_location_y();
			u.Speed=rdr.read_unit_speed();
			u.Direction=rdr.read_unit_face();
			u.AimingFace=rdr.read_unit_aiming_face();
			wtr.write_unit_id(u.ID);
			wtr.write_uint8(u.Type);
			wtr.write_uint8(u.Camps);
			wtr.write_uint8(u.Level);
			wtr.write_uint8(u.Ranks);
			wtr.write_uint16(u.HP);
			wtr.write_uint16(u.X);
			wtr.write_uint16(u.Y);
			wtr.write_uint16(u.Speed);
			wtr.write_uint16(u.Direction);
			wtr.write_uint16(u.AimingFace);
		})
	}
	ph2(utils.CMD_create_unit);
	ph1(uint16(wtr.pos)-2);
	return res;
}