package command

import (
	"github.com/pkg/errors"
)

type IUnitCreateDataSetter interface {
	IUnitMovementDataSetter;
	SetCamps(uint8);
	SetAttackPower(uint16);
	SetHP(uint16);
	SetStat(uint16);
}

func CreateUnit(data []byte,build_unit func(uint16)IUnitCreateDataSetter)(error){
	r:=&packet_reader{
		data:data,
		pos:0,
	}
	count:=(int)(r.read_unit_count());
	for i:=0;i<count;i++{
		if u:=build_unit(r.read_unit_id());u!=nil{
			u.SetCamps(
				r.read_unit_camps(),
			);
			u.SetAttackPower(
				r.read_unit_attk_power(),
			);
			u.SetHP(
				r.read_unit_hp(),
			);
			u.SetLocation(
				r.read_unit_location_x(),
				r.read_unit_location_y(),
			);
			u.SetMovement(
				r.read_unit_speed(),
				r.read_unit_face(),
				r.read_unit_aiming_face(),
			);
			u.SetStat(
				r.read_unit_stat(),
			);
		}else{
			return errors.New("packet error");
		}
	}
	return nil;
}