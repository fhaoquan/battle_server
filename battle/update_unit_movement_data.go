package battle

import (
	"github.com/pkg/errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)UpdateUnitMovement(pkt []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	r:=&packet_decoder{
		data:pkt,
		pos:0,
	}
	count:=(int)(r.read_unit_count());
	for i:=0;i<count;i++{
		if u:=context.FindUnit(r.read_unit_id());u!=nil{
			u.X=r.read_unit_location_x();
			u.Y=r.read_unit_location_y();
			u.Speed=r.read_unit_speed();
			u.Direction=r.read_unit_face();
			u.AimingFace=r.read_unit_aiming_face();
		}else{
			return errors.New("packet error");
		}
	}
	return nil;
}
