package player

import (
	"../utils"
	"net"
	"fmt"
	"errors"
	"github.com/sirupsen/logrus"
)
type Error interface {
	error;
}
type BattlePlayer struct {
	uid uint32;
	name string;
	login_times int;
	kcp_conn net.Conn;
	udp_conn net.PacketConn;
	udp_addr net.Addr;
	kcp_out chan utils.IKcpRequest;
	udp_out chan utils.IUdpRequest;
	close_sig chan interface{};
	kcp_req_pool *utils.MemoryPool;
}

func (me *BattlePlayer)Connected()bool {
	return me.login_times>0;
}
func (me *BattlePlayer)UID()(uint32){
	return me.uid;
}
func (me *BattlePlayer)New(id uint32,name string)(*BattlePlayer){
	me.uid=id;
	me.name=name;
	me.login_times=0;
	return me;
}
func (me *BattlePlayer)GetUdpMsgChan()(chan utils.IUdpRequest){
	return me.udp_out;
}
func (me *BattlePlayer)GetKcpMsgChan()(chan utils.IKcpRequest){
	return me.kcp_out;
}
func (me *BattlePlayer)ListenUDP()(err error){
	return nil;
}
func (me *BattlePlayer)SendKCP(b []byte)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	_,err=me.kcp_conn.Write(b);
	return ;
}
func (me *BattlePlayer)SendUDP(b []byte)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	_,err=me.udp_conn.WriteTo(b,me.udp_addr);
	return ;
}
func (me *BattlePlayer)RecvKCP(empty utils.IKcpRequest)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	if e:=empty.ReadAt(me.kcp_conn);e!=nil{
		return e;
	}
	if empty.GetUID()!=me.uid{
		return errors.New("error uid");
	}
	return nil;
}
func (me *BattlePlayer)RecvUDP(empty utils.IUdpRequest)(err error){
	defer func(){
		if e:=recover();e!=nil{
			err=errors.New(fmt.Sprint(e));
		}
	}()
	if e:=empty.ReadAt(me.udp_conn);e!=nil{
		return e;
	}
	if empty.GetUID()!=me.uid{
		return errors.New("error uid");
	}
	me.udp_addr=empty.GetAdr();
	return nil;
}
func (me *BattlePlayer)ResetKCP(conn net.Conn){
	me.kcp_conn=conn;
	me.login_times++;
	go me.kcp_recv_proc(me.uid,conn);
}
func (me *BattlePlayer)kcp_recv_proc(uid uint32,conn net.Conn){
	f:=func(empty *KcpReq)(err error){
		defer func(){
			if e:=recover();e!=nil{
				err=errors.New(fmt.Sprint(e));
			}
		}()
		if e:=empty.ReadAt(conn);e!=nil{
			return e;
		}
		if empty.GetUID()!=uid{
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
		err_times:=0;
		for{
			select {
			case _,ok:=<-me.close_sig:
				if !ok {
					return nil;
				}
			default:
				r:=me.kcp_req_pool.Pop().(*KcpReq);
				switch e:=f(r);e.(type){
				case nil:
					me.kcp_out<-r;
				case net.Error:
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
	}
}