package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)UnitAttackStart(data []byte)(i interface{}){
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()

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
	cnt:=rdr.read_uint8();
	out:=0;
	ph1:=wtr.get_uint16_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	ph3:=wtr.get_uint08_placeholder();
	for i:=0;i<(int)(cnt);i++{
		uid:=rdr.read_uint16();
		tid:=rdr.read_uint16();
		u:=context.FindUnit(uid);
		if u==nil{
			continue;
		}
		if u.Death(){
			continue;
		}
		out++;
		wtr.write_uint16(uid);
		wtr.write_uint16(tid);
	}
	ph3(byte(out));
	ph2(utils.CMD_attack_start);
	ph1(uint16(wtr.pos)-2);
	return res;
}