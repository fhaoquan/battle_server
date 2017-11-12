package battle

import (
	"github.com/sirupsen/logrus"
	"fmt"
	"runtime/debug"
)

func (context *Battle)ComputeScore(who uint32)(res uint32){
	res=0;
	defer func(){
		if e:=recover();e!=nil{
			logrus.Error(e);
			logrus.Error(fmt.Sprintf("%s",debug.Stack()));
			res=0;
		}
	}()
	uid:=uint16(0);
	for e:=context.main_base_list.Front();e!=nil;e=e.Next(){
		uid=e.Value.(uint16);
		context.FindUnitDo(uint16(uid), func(u *Unit) {
			if u.Killer==who&&u.Death(){
				res+=uint32(u.Score);
			}
		})
	}
	return ;
}
