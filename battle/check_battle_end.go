package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)CheckBattleEnd()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	end:=false;
	uid:=uint16(0);
	for e:=context.main_base_list.Front();e!=nil&&!end;e=e.Next(){
		uid=e.Value.(uint16);
		context.FindUnitDo(uint16(uid), func(u *Unit) {
			if u.Death(){
				end=true;
			}
		})
	}
	if !end{
		return nil;
	}
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.UID=0;
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	wtr.write_uint16(uint16(6));
	wtr.write_uint8(utils.CMD_battle_end)
	wtr.write_uint16(uid);
	return res;
}
