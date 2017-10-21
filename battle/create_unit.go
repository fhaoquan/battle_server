package battle

import (
	"fmt"
	"errors"
	"../utils"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)CreateUnit(who uint32,data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	logrus.Error("in CreateUnit");
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
	ph1:=wtr.get_uint16_placeholder();
	ph2:=wtr.get_uint08_placeholder();
	count:=rdr.read_unit_count();
	wtr.write_unit_count(count);
	for i:=0;i<(int)(count);i++{
		context.CreateUnitDo(func(u *Unit) {
			u.Owner=who;
			u.Type=rdr.read_uint32();
			u.Card=rdr.read_uint32();
			u.Level=rdr.read_uint8();
			u.Ranks=rdr.read_uint8();
			u.HP=rdr.read_unit_hp();
			u.X=rdr.read_unit_location_x();
			u.Y=rdr.read_unit_location_y();
			u.WriteToBuf(wtr);
		})
	}
	ph2(utils.CMD_create_unit);
	ph1(uint16(wtr.pos)-2);
	return res;
}