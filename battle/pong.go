package battle

import (
	"fmt"
	"../utils"
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)Pong(who uint32,data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	res:=context.kcp_res_pool.Pop().(*utils.KcpRes);
	res.Broadcast=false;
	res.UID=who;
	wtr:=&packet_encoder{
		res.BDY,
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	wtr.write_bytes(data);
	ph1(utils.CMD_pingpong);
	ph0(uint16(wtr.pos)-2);
	return res;
}
