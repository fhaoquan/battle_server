package room

import (
	"net"
	"fmt"
	"errors"
	"github.com/sirupsen/logrus"
	"../utils"
)
func (me *Room1v1)kcp_recv_proc(session *kcp_session){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		err_times:=0;
		pool:=utils.NewMemoryPool(16, func(impl utils.ICachedData) utils.ICachedData {
			return &KcpReq{
				impl,0,0,0,make([]byte,utils.MaxPktSize),
			}
		})
		for{
			select {
			case _,ok:=<-me.close_sig:
				if !ok {
					return nil;
				}
			default:
				r:=pool.Pop().(*KcpReq);
				switch e:=session.Recv(r);e.(type){
				case nil:
					err_times=0;
					me.kcp_chan<-r;
				case net.Error:
					logrus.Error(e);
					r.Return();
					if e.(net.Error).Temporary(){
						if err_times++;err_times>5{
							return errors.New("5 times temporary error!");
						}
					}else{
						return e;
					}
				default:
					r.Return();
					return e;
				}
			}
		}
	}();
	session.Close();
	if e!=nil{
		logrus.Error(e);
		return;
	}
}
