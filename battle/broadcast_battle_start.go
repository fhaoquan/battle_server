package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleStart(uid uint32)(i interface{}){
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res.UID=uid;
	if uid==0{
		res.Broadcast=false;
	}else{
		res.Broadcast=true;
	}
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	wtr.write_uint16(2);
	wtr.write_uint8(utils.CMD_battle_start);
	wtr.write_uint8(0);
	return res;
}