package battle

import (
	"fmt"
	"errors"
	"../utils"
)

func (context *Battle)UnitAttackStart(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	f:=wtr.get_uint16_placeholder();
	wtr.write_bytes(data);
	res.LEN=uint16(wtr.pos);
	f(res.LEN-2);
	return res;
}