package battle

import (
	"fmt"
	"errors"
	"time"
	"../utils"
)

func (context *Battle)BroadcastBattleWaitingStart(t time.Time)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	now:=time.Now();
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.UID=0;
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	res.Broadcast=true;
	wtr.write_uint16(9);
	wtr.write_uint8(utils.CMD_battle_wating_start);
	if t.After(now){
		wtr.write_uint64(uint64(t.Sub(t).Nanoseconds()));
	}else{
		wtr.write_uint64(0);
	}

	return res;

}
