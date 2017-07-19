package command

import (
	"github.com/pkg/errors"
)

type IUnitMovementDataSetter interface {
	SetLocation(x uint16,y uint16);
	SetMovement(speed uint16,face uint16,aiming_face uint16);
}
func UpdateUnitMovement(data []byte,finder func(uint16)IUnitMovementDataSetter)(error){
	r:=&packet_decoder{
		data:data,
		pos:0,
	}
	count:=(int)(r.read_unit_count());
	for i:=0;i<count;i++{
		if u:=finder(r.read_unit_id());u!=nil{
			u.SetLocation(
				r.read_unit_location_x(),
				r.read_unit_location_y());
			u.SetMovement(
				r.read_unit_speed(),
				r.read_unit_face(),
				r.read_unit_aiming_face());
		}else{
			return errors.New("packet error");
		}
	}
	return nil;
}