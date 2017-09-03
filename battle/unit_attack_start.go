package battle

import (
	"fmt"
	"errors"
)

func (context *Battle)UnitAttackStart(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=context.kcp_res_pool.Pop().(*kcp_response);
	res.broadcast=true;
	wtr:=&packet_encoder{
		res.bdy,
		0,
	}
	f:=wtr.get_uint16_placeholder();
	wtr.write_bytes(data);
	res.len=uint16(wtr.pos);
	f(res.len-2);
	return res;
}