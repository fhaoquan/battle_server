package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)
func (context *Battle)BroadcastBattleAll(uid uint32)(i interface{}){
	res:=context.kcp_res_pool.Pop().(utils.I_RES);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res.SetUID(uid);
	if uid==0{
		res.SetBroadcast(false);
	}else{
		res.SetBroadcast(true);
	}

	wtr:=&packet_encoder{
		res.GetWriteBuffer(),
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=0;
	for e:=context.living_units.Front();e!=nil;e=e.Next(){
		uid:=e.Value.(uint16);
		context.FindUnitDo(uint16(uid), func(u *Unit) {
			u.WriteToBuf(wtr);
			count++;
		})
	}
	ph2(uint8(count));
	ph1(utils.CMD_battle_all);
	ph0(uint16(wtr.pos)-2);
	return res;
}
