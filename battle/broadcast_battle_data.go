package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleMovementData(receiver uint32, owner_filter uint32)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res:=context.udp_res_pool.Pop().(*utils.UdpRes);
	//res.UID=0;
	res.Broadcast=false;
	res.UID=receiver;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=0;

	for e:=context.living_units.Front();e!=nil;e=e.Next(){
		uid:=e.Value.(uint16);
		context.FindUnitDo(uint16(uid), func(u *Unit) {
			if(u.Owner==owner_filter){
				wtr.write_unit_id(u.ID).
					write_unit_x(u.X).
					write_unit_y(u.Y).
					write_unit_speed(u.Speed).
					write_unit_face(u.Direction).
					write_unit_aiming_face(u.AimingFace);
				count++;
			}
		})
	}
	ph2(uint8(count));
	ph1(utils.CMD_unit_movment);
	ph0(uint16(wtr.pos)-2);
	if count>0{
		return res;
	}else{
		res.Return();
		return nil;
	}

}
