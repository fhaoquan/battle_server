package gateway

import (
	"net"
	"time"
	"../utils"
	"github.com/sirupsen/logrus"
	"sync"
	"encoding/binary"
	"fmt"
	"runtime/debug"
	"errors"
)
type pending_packet struct {
	utils.ICachedData;
	addr			net.Addr;
	len				int;
	data			[]byte;
}
type kcp_wrapper struct {
	m					*sync.Mutex;
	kcp					*KCP;
}
func (this *kcp_wrapper)conv()(uint32){
	return this.kcp.conv;
}
func (this *kcp_wrapper)send(data []byte)(int){
	defer this.m.Unlock();
	this.m.Lock();
	this.kcp.Send(data);
	this.kcp.flush(false);
	return 0;
}
func (this *kcp_wrapper)read(udp_data []byte,kcp_data []byte)(bool,error){
	defer this.m.Unlock();
	this.m.Lock();
	if this.kcp.Input(udp_data,true,true)<0{
		return false,errors.New("kcp.Input error")
	}
	if this.kcp.PeekSize()<=0{
		return false,nil;
	}
	if this.kcp.Recv(kcp_data)>=0{
		return true,nil;
	}else{
		return false,errors.New("kcp.Recv error");
	}
}
type session struct{
	wait_close			sync.WaitGroup;
	once_close 			sync.Once;
	once_start 			sync.Once;
	conn				*net.UDPConn;
	kcp					*kcp_wrapper;
	uid					uint32;
	rid					uint32;
	chan_die_signal		chan interface{}
	remote_addr			net.Addr;
}
func (me *session)Start(k chan utils.I_REQ){
	me.start(k);
}
func (me *session)RoomID()(uint32){
	return me.rid;
}
func (me *session)UserID()(uint32){
	return me.uid;
}
func (me *session)Close(){
	me.close(false);
}
func (me *session)KcpSend(data []byte)(bool){
	_,e:=me.send_kcp(data);
	return e!=nil;
}
func (me *session)UdpSend(data []byte)(bool){
	_,e:=me.send_udp(data);
	return e!=nil;
}
func (me *session)go_session_kernel_proc(f func(i ...interface{}),i ...interface{}){
	me.wait_close.Add(1);
	go func(){
		defer func() {
			if e:=recover();e!=nil{
				logrus.Error(e);
				logrus.Error(fmt.Sprintf("%s",debug.Stack()));
			}
			me.wait_close.Done();
			me.close(false);
		}()
		f(i...);
	}();
}
func (me *session)close(wait bool){
	me.once_close.Do(func() {
		close(me.chan_die_signal);
		receiver:=GetReceiver(me.rid,me.uid);
		if receiver!=nil {
			receiver.DelSession(me);
		}
	})
	if wait{
		me.wait_close.Wait();
	}
}
func (me *session)start(k chan utils.I_REQ){
	me.once_start.Do(func() {
		me.go_session_kernel_proc(me.recv_proc,k);
	})
}
func (me *session)send_kcp(data []byte)(int,error){
	select {
	case <-me.chan_die_signal:
		return 0,errors.New("session closed");
	default:
		return me.kcp.send(data),nil;
	}
}
func (me *session)send_udp(data []byte)(int,error){
	select {
	case <-me.chan_die_signal:
		return 0,errors.New("session closed");
	default:
		return me.conn.WriteTo(data,me.remote_addr);
	}
}
func (me *session)recv_proc(i ...interface{}){
	receiver:=i[0].(chan utils.I_REQ);
	kcp_pool:=utils.NewKcpReqPool(16);
	udp_pool:=utils.NewUdpReqPool(16);
	for{
		udp_pkt:=udp_pool.Pop().(utils.I_REQ)
		me.conn.SetReadDeadline(time.Now().Add(time.Second*5));
		l,a,e:=me.conn.ReadFrom(udp_pkt.GetALL());
		if e!=nil{
			udp_pkt.Return();
			logrus.Error(e);
			return;
		}
		h:=binary.LittleEndian.Uint32(udp_pkt.GetALL()[0:4]);
		switch h{
		case me.kcp.conv():
			kcp_pkt:=kcp_pool.Pop().(utils.I_REQ);
			y,e:=me.kcp.read(udp_pkt.GetALL()[:l],kcp_pkt.GetALL());
			udp_pkt.Return();
			if e!=nil{
				kcp_pkt.Return();
				logrus.Error("kcp protocol error");
				return ;
			}
			if !y{
				kcp_pkt.Return();
				continue;
			}
			if !kcp_pkt.Check(){
				kcp_pkt.Return();
				logrus.Error("request protocol error");
				return ;
			}
			me.remote_addr=a;
			select {
			case receiver<-kcp_pkt:
				continue;
			case <-me.chan_die_signal:
				kcp_pkt.Return();
				return ;
			}
		case utils.UdpPktHeader:
			if !udp_pkt.Check(){
				udp_pkt.Return();
				logrus.Error("request protocol error");
				return ;
			}
			me.remote_addr=a;
			select {
			case receiver<-udp_pkt:
				continue;
			case <-me.chan_die_signal:
				udp_pkt.Return();
				return ;
			}
		default:
			logrus.Error("unknown packet type=",h);
			udp_pkt.Return();
		}

	}
}
func new_session(rid uint32,uid uint32)(*session,error){
	a,e:= net.ResolveUDPAddr("udp", ":0");
	if e!=nil{
		return nil,e;
	}
	c,e:=net.ListenUDP("udp",a);
	session:=new(session);
	session.remote_addr=a;
	kcp:=NewKCP(uid, func(buf []byte, size int) {
		c.WriteTo(buf[:size],session.remote_addr)
	});
	kcp.stream=1;
	kcp.mtu=1350;
	kcp.NoDelay(1,5,2,1);
	kcp.WndSize(32,32);

	session.kcp=&kcp_wrapper{new(sync.Mutex),kcp};
	session.conn=c;
	session.rid=rid;
	session.uid=uid;
	session.chan_die_signal=make(chan interface{},1);
	return session,nil;
}
