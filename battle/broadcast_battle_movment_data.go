package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleMovementData(receiver uint32, owner_filter uint32,status uint16,remaining uint16)(i interface{}){
	res:=context.udp_res_pool.Pop().(utils.I_RES);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()

	//res.UID=0;
	res.SetBroadcast(false);
	res.SetUID(receiver);
	wtr:=&packet_encoder{
		res.GetWriteBuffer(),
		0,
	}
	//wtr.write_uint32(utils.UdpPktHeader);
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	wtr.write_uint16(status).
		write_uint16(remaining);
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
					write_unit_aiming_face(u.AimingFace).
					write_uint16(u.Status);
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
