package battle

import (
	"../utils"
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime/debug"
	"errors"
)

func (context *Battle)BroadcastBattleRemainingTime(status uint16,remaining uint16)(i interface{}){
	res:=context.kcp_res_pool.Pop().(utils.I_RES);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res.SetUID(0);
	res.SetBroadcast(true);
	wtr:=&packet_encoder{
		res.GetWriteBuffer(),
		0,
	}
	wtr.write_uint16(5);
	wtr.write_uint8(utils.CMD_battle_remaining_time);
	wtr.write_uint16(status);
	wtr.write_uint16(remaining);
	return res;
}
