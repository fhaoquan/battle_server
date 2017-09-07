package battle

import (
	"fmt"
	"errors"
	"../utils"
)
func (context *Battle)BroadcastBattleAll()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=context.udp_res_pool.Pop().(*utils.UdpRes);
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	count:=0;
	context.ForEachUnitDo(func(u *Unit)bool{
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
	res.LEN=uint16(wtr.pos);
	ph0(res.LEN-2);
	return res;
}
