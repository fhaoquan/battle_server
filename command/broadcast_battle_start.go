package command

import (
	"fmt"
	"errors"
	"time"
)

func (cmd *CommandContext)BroadcastBattleStart()(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=cmd.udp_res_pool.Pop().(*udp_response);
	wtr:=&packet_encoder{
		res.bdy,
		0,
	}
	wtr.write_uint16(8);
	wtr.write_uint64((uint64)(time.Now().Add(time.Second*5).UTC().Unix()));
	return res;
}
