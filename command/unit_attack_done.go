package command

import (
	"../utils"
	"../sessions/packet"
	"errors"
	"fmt"
)
func (cmd *Commamd)UnitAttackDone(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
		}
	}()
	w:=&packet_encoder{
		make([]byte,utils.MaxPktSize),
		0,
	}
	r:=&packet_decoder{
		data:data,
		pos:0,
	}
	power:=(r.read_unit_attack_power())
	count:=(int)(r.read_unit_count());
	w.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		id:=r.read_unit_id();
		w.write_unit_id(id);
		u2:=cmd.base_room.GetBattle().FindUnit(id);
		if u2!=nil{
			if u2.HP>power{
				u2.HP-=power;
			}else{
				u2.HP=0;
			}
			w.write_uint16(u2.HP);
		}else{
			w.write_uint16(0);
		}
	}
	return packet.IKcpResponse(&struct {
		broadcast bool;
		uid uint32;
		bdy []byte;
	}{
		true,
		0,
		w.data[:w.pos],
	});
}
