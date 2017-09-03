package command

import (
	"errors"
	"fmt"
)
func (cmd *CommandContext)each_unit_attack_done(rdr *packet_decoder,wtr *packet_encoder)(){
	power:=(rdr.read_unit_attack_power())
	count:=(int)(rdr.read_unit_count());
	wtr.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		id:=rdr.read_unit_id();
		wtr.write_unit_id(id);
		u2:=cmd.base_room.GetBattle().FindUnit(id);
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
func (cmd *CommandContext)UnitAttackDone(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	res:=cmd.kcp_res_pool.Pop().(*kcp_response);
	res.broadcast=true;
	wtr:=&packet_encoder{
		res.bdy,
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
		cmd.each_unit_attack_done(rdr,wtr);
	}
	res.len=(uint16)(wtr.pos);
	f(res.len-2);
	return res;
}
