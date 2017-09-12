package kcp_server

import (
	"../../utils"
	"net"
	"github.com/sirupsen/logrus"
	"time"
	"sync"
)
var (
	xmitBuf sync.Pool;
)

func init() {
	xmitBuf.New = func() interface{} {
		return make([]byte,utils.MaxPktSize)
	}
}
type KcpSyncWrapper struct {
	mu			sync.Mutex;
	kcp			*KCP;
}
func (me *KcpSyncWrapper)use_kcp_in_mux(do func(k *KCP)){
	defer me.mu.Unlock();
	me.mu.Lock();
	do(me.kcp);
}
func (me *KcpSyncWrapper)unsafe_use_kcp()*KCP {
	return me.kcp;
}

type KcpSession struct {
	wait_close	sync.WaitGroup;
	once_close 	sync.Once;
	once_start 	sync.Once;
	getway		*KcpGateway;
	conn		net.PacketConn
	kcp			*KcpSyncWrapper;
	chPending	chan *pending_packet;
	chDie		chan interface{};
	ChRecv		chan utils.IKcpRequest;
	RemoteAddr	net.Addr;
}
func (session *KcpSession)go_session_kernel_proc(f func(a interface{}),a interface{}){
	session.wait_close.Add(1);
	go func(){
		defer func() {
			session.wait_close.Done();
			if e:=recover();e!=nil{
				logrus.Error(e);
				session.Close(false);
			}
		}()
		f(a);
	}();
}
func (session *KcpSession)start_recv_proc(){
	session.go_session_kernel_proc(func(a interface{}){
		pool:=utils.NewMemoryPool(16, func(impl utils.ICachedData) utils.ICachedData {
			return &utils.KcpReq{
				impl,make([]byte,utils.MaxPktSize),
			}
		})
		for{
			select {
			case d,ok:=<-session.chPending:
				if ok{
					p:=pool.Pop().(*utils.KcpReq);
					session.kcp.use_kcp_in_mux(func(k *KCP) {
						k.Input(d.data[:d.len],true,true);
						if n:=k.PeekSize();n>0{
							k.Recv(p.GetALL());
							if p.Check() {
								session.ChRecv<-p;
							}else{
								logrus.Error("recved fail :",p.GetALL());
								p.Return();
							}
						}else{
							p.Return();
						}
					})
					d.Return();
				}
			case <-session.chDie:
				return ;
			}
		}
	},nil);
}
func (session *KcpSession)start_updt_proc(){
	session.go_session_kernel_proc(func(a interface{}){
		d:=time.Millisecond;
		for{
			select {
			default:
				time.Sleep(d);
				session.kcp.use_kcp_in_mux(func(k *KCP) {
					k.flush(false);
					d=time.Duration(k.interval) * time.Millisecond
				})

			case <-session.chDie:
				return ;
			}
		}
	},nil);
}
func (session *KcpSession)kcp_send_callback(buf []byte, size int){
	session.conn.WriteTo(buf[:size],session.RemoteAddr);
}
func (session *KcpSession)Send(data []byte){
	session.kcp.use_kcp_in_mux(func(k *KCP) {
		k.Send(data);
		//k.flush(false);
	})
}
func (session *KcpSession)Start(){
	session.once_start.Do(func() {
		logrus.Error("session :",session.RemoteAddr," started");
		session.start_updt_proc();
		session.start_recv_proc();
	})
}
func (session *KcpSession)Close(wait bool){
	session.once_close.Do(func() {
		session.getway.del_session(session);
		close(session.chDie);
		logrus.Error("session :",session.RemoteAddr," closed");
	})
	if wait{
		session.wait_close.Wait();
	}
}
func (session *KcpSession)congestion()bool{
	return !(len(session.chPending)<cap(session.chPending));
}
func (session *KcpSession) GetConv() uint32 {
	return session.kcp.unsafe_use_kcp().conv;
}
func (session *KcpSession) LocalAddr() net.Addr {
	return session.conn.LocalAddr();
}
func new_kcp_session(getway *KcpGateway,conv uint32,conn net.PacketConn,remote net.Addr)(*KcpSession){
	s:=new(KcpSession);
	s.getway=getway;
	s.RemoteAddr=remote;
	s.conn=conn;
	s.kcp=new(KcpSyncWrapper);
	s.kcp.kcp=NewKCP(conv,s.kcp_send_callback);
	s.kcp.unsafe_use_kcp().stream=1;
	s.kcp.unsafe_use_kcp().mtu=1350;
	s.kcp.unsafe_use_kcp().NoDelay(1,5,2,1);
	s.kcp.unsafe_use_kcp().WndSize(32,32);
	s.chPending=make(chan *pending_packet,16);
	s.chDie=make(chan interface{},1);
	s.ChRecv=make(chan utils.IKcpRequest,16);
	return s;
}
