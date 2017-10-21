package battle

import (
	"errors"
	"fmt"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)
func (context *Battle)each_unit_attack_done(rdr *packet_decoder,wtr *packet_encoder)(){
	effect:=rdr.read_uint32();
	wtr.write_uint32(effect);
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
				logrus.Error("unit hp==",u2.HP," id==",u2.ID);
			}else{
				logrus.Error("unit hp==0 id==",u2.ID);
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
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	logrus.Error("UnitAttackDone 1");
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	logrus.Error("UnitAttackDone 2");
	res.UID=0;
	res.Broadcast=true;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	rdr:=&packet_decoder{
		data:data,
		pos:0,
	}
	ph1:=wtr.get_uint16_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=(int)(rdr.read_unit_count());
	wtr.write_unit_count((uint8)(count));
	for i:=0;i<count;i++{
		context.each_unit_attack_done(rdr,wtr);
	}
	ph2(utils.CMD_attack_done);
	ph1((uint16)(wtr.pos)-2);
	logrus.Error("UnitAttackDone 3");
	return res;
}
