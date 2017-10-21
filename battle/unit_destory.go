package battle

import (
	"fmt"
	"errors"
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func (context *Battle)UnitDestory(data []byte)(i interface{}){
	defer func(){
		if e:=recover();e!=nil{
			i=errors.New(fmt.Sprint(e));
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
		}
	}()
	rdr:=&packet_decoder{
		data:data,
		pos:0,
	}
	cnt:=rdr.read_uint8();
	for i:=0;i<(int)(cnt);i++{
		uid:=rdr.read_uint16();
		for e:=context.living_units.Front();e!=nil;e=e.Next(){
			if ((e.Value))==(uid){
				context.living_units.Remove(e);
				break;
			}
		}
	}
	return nil;
}
