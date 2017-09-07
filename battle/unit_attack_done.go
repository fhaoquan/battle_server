package battle

import (
	"errors"
	"fmt"
	"../utils"
)
func (context *Battle)each_unit_attack_done(rdr *packet_decoder,wtr *packet_encoder)(){
	power:=(rdr.read_unit_attack_power())
	count:=(int)(rdr.read_unit_count());
	wtr.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		id:=rdr.read_unit_id();
		wtr.write_unit_id(id);
		u2:=context.FindUnit(id);
		if u2!=nil{
			if u2.HP>power{
				u2.HP-=power;
			}else{
				u2.HP=0;
			}
			wtr.write_uint16(u2.HP);
		}else{
			wtr.write_uint16(0);
		}
	}
}
func (context *Battle)UnitAttackDone(data []byte)(i interface{}){
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
	rdr:=&packet_decoder{
		data:data,
		pos:0,
	}
	f:=wtr.get_uint16_placeholder();
	count:=(int)(rdr.read_unit_count());
	wtr.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		context.each_unit_attack_done(rdr,wtr);
	}
	res.LEN=(uint16)(wtr.pos);
	f(res.LEN-2);
	return res;
}
