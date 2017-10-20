package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)BroadcastBattleStart()(i interface{}){
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
	res.Broadcast=true;
	wtr.write_uint16(2);
	wtr.write_uint8(utils.CMD_battle_start);
	wtr.write_uint8(0);
	return res;

}