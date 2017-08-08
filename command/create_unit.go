package command

import (
	"../utils"
	"../battle"
	"../sessions/packet"
	"fmt"
	"errors"
)

func (cmd *Commamd)CreateUnit(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	w:=&packet_encoder{
		make([]byte,utils.MaxPktSize),
		0,
	}
	r:=&packet_decoder{
		data:data,
		pos:0,
	}
	count:=r.read_unit_count();
	w.write_unit_count(count);
	for i:=0;i<(int)(count);i++{
		cmd.base_room.GetBattle().CreateUnitDo(r.read_unit_id(), func(u *battle.Unit) {
			u.Camps=r.read_unit_camps();
			u.HP=r.read_unit_hp();
			u.X=r.read_unit_location_x();
			u.Y=r.read_unit_location_y();
			u.Speed=r.read_unit_speed();
			u.Direction=r.read_unit_face();
			u.AimingFace=r.read_unit_aiming_face();
			w.write_unit_id(u.ID)
			w.write_uint8(u.Camps)
			w.write_uint16(u.HP);
			w.write_uint16(u.X);
			w.write_uint16(u.Y);
			w.write_uint16(u.Speed);
			w.write_uint16(u.Direction);
			w.write_uint16(u.AimingFace);
		})
	}
	return packet.IKcpResponse(&struct {
		broadcast bool;
		uid uint32;
		bdy []byte;
	}{
		true,
		0,
		w.data[:w.pos],
	});
}