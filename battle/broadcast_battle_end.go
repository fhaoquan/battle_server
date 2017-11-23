package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleEnd(winner uint32)(i interface{}){
	res:=(utils.I_RES)(nil);
	defer func(){
		if e:=recover();e!=nil{
			if res!=nil {
				res.Return();
			}
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res=context.kcp_res_pool.Pop().(utils.I_RES);
	res.SetUID(0);
	res.SetBroadcast(true);
	wtr:=&packet_encoder{
		res.GetWriteBuffer(),
		0,
	}
	wtr.write_uint16(uint16(5));
	wtr.write_uint8(utils.CMD_battle_end)
	wtr.write_uint32(winner);
	return res;
}
