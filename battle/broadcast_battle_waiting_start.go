package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleWaitingStart()(i interface{}){
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
	wtr.write_uint16(2);
	wtr.write_uint8(utils.CMD_battle_wating_start);
	wtr.write_uint8(0);
	return res;

}
