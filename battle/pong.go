package battle

import (
	"fmt"
	"../utils"
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)Pong(who uint32,data []byte)(i interface{}){
	res:=context.kcp_res_pool.Pop().(utils.I_RES);
	defer func(){
		if e:=recover();e!=nil{
			res.Return();
			i=&BattlePanicError{errors.New(fmt.Sprint(e))};
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()

	res.SetBroadcast(false);
	res.SetUID(who);
	wtr:=&packet_encoder{
		res.GetWriteBuffer(),
		0,
	}
	ph0:=wtr.get_uint16_placeholder();
	ph1:=wtr.get_uint08_placeholder();
	wtr.write_bytes(data);
	ph1(utils.CMD_pingpong);
	ph0(uint16(wtr.pos)-2);
	return res;
}
