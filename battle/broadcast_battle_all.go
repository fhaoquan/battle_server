package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)
func (context *Battle)BroadcastBattleAll()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.UID=0;
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=0;
	context.ForEachUnitDo(func(u *Unit)bool{
		u.WriteToBuf(wtr);
		count++;
		return true;
	})
	ph2(uint8(count));
	ph1(utils.CMD_battle_all);
	ph0(uint16(wtr.pos)-2);
	return res;
}
