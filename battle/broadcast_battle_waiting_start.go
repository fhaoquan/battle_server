package battle

import (
	"fmt"
	"errors"
	"../utils"
)

func (context *Battle)BroadcastBattleWaitingStart()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
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
	wtr.write_uint8(utils.CMD_battle_wating_start);
	wtr.write_uint8(0);
	return res;

}
