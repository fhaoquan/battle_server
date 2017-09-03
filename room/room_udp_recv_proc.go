package room

import (
	"../utils"
	"fmt"
	"errors"
	"net"
	"github.com/sirupsen/logrus"
)

func (me *Room1v1)udp_recv_proc(conn net.PacketConn){
	defer func(){
		me.wait.Done();
	}()
	me.wait.Add(1);
	f:= func(empty utils.IUdpRequest)(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		if e:=empty.ReadAt(conn);e!=nil{
			return e;
		}
		switch empty.GetUID() {
		case me.p1.uid:
			me.p1.udp_session.udp_addr=empty.GetAdr();
		case me.p2.uid:
			me.p2.udp_session.udp_addr=empty.GetAdr();
		default:
			return errors.New("error uid");
		}
		return nil;
	}
	e:=func()(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()

		pool:=utils.NewMemoryPool(16, func(impl utils.ICachedData) utils.ICachedData {
			return &UdpReq{
				impl,&net.UDPAddr{},0,0,0,make([]byte,utils.MaxPktSize),
			}
		})
		err_times:=0;
		for{
			select {
			case _,ok:=<-me.close_sig:
				if !ok {
					return nil;
				}
			default:
				r:=pool.Pop().(*UdpReq);
				switch e:=f(r);e.(type){
				case nil:
					err_times=0;
					me.udp_chan<-r;
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
	conn.Close();
	if e!=nil{
		logrus.Error(e);
		logrus.Error("udp recv proc exited !!")
		return;
	}
}
