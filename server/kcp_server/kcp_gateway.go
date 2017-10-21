package kcp_server

import (
	"../../utils"
	"net"
	"github.com/sirupsen/logrus"
	"sync"
	"encoding/binary"
	"time"
	"golang.org/x/net/ipv4"
	"runtime/debug"
	"fmt"
)
type pending_connection struct{
	uid				uint32;
	rid				uint32;
}
type pending_packet struct {
	utils.ICachedData;
	addr			net.Addr;
	len				int;
	data			[]byte;
}
func (this *pending_packet)get_uid()uint32{
	return binary.LittleEndian.Uint32(this.data[2:]);
}
func (this *pending_packet)get_rid()uint32{
	return binary.LittleEndian.Uint32(this.data[6:]);
}
func (this *pending_packet)check_handshake()bool{
	if binary.LittleEndian.Uint16(this.data) !=54321{
		return false;
	}
	if this.len!=10{
		return false;
	}
	return true;
}
type KcpGateway struct{
	once_start			sync.Once;
	conn				*net.UDPConn;
	sessions			sync.Map;
	on_session 			func(uid,rid uint32,session *KcpSession);
}

func (this *KcpGateway)go_gateway_kernel_proc(f func(a interface{}),a interface{}){
	go func(){
		defer func() {
			if e:=recover();e!=nil{
				logrus.Error(e);
				logrus.Error(fmt.Sprintf("%s",debug.Stack()));
			}
		}()
		f(a);
	}();
}

func (this *KcpGateway)start_main_proc(){
	ch_handshake:=make(chan *pending_packet,64);
	this.go_gateway_kernel_proc(func(a interface{}) {
		evnt:=a.(chan *pending_packet);
		for{
			select {
			case p,ok:=<-evnt:
				if ok{
					if p.check_handshake(){
						c:=&pending_connection{
							p.get_uid(),
							p.get_rid(),};
						this.sessions.Store(p.addr.String(),c);
						this.conn.WriteTo(p.data[:p.len],p.addr);
						go func(key interface{}){
							time.Sleep(time.Minute);
							if s,ok:=this.sessions.Load(key);ok{
								switch s.(type){
								case (*pending_connection):
									this.sessions.Delete(key);
								}
							}
						}(p.addr.String())
					}else{
						p.Return();
					}
				}
			}
		}
	},ch_handshake);
	this.go_gateway_kernel_proc(func(a interface{}) {
		evnt:=a.(chan *pending_packet);
		pool:=utils.NewMemoryPool(64*1000, func(impl utils.ICachedData) utils.ICachedData {
			return &pending_packet{
				impl,nil,0,make([]byte,1500),
			}
		})
		e:=error(nil);
		for{
			p:=pool.Pop().(*pending_packet);
			p.len,p.addr,e=this.conn.ReadFrom(p.data);
			if p.len<4{
				continue;
			}
			if(e!=nil){
				logrus.Error(e);
				p.Return();
				continue;
			}
			if s,ok:=this.sessions.Load(p.addr.String());ok{
				switch s.(type){
				case *pending_connection:
					if this.on_session!=nil{
						conv:=binary.LittleEndian.Uint32(p.data);
						new_session:=new_kcp_session(this,conv,this.conn,p.addr);
						new_session.Start();
						this.add_session(new_session);
						new_session.chPending<-p;
						go this.on_session(s.(*pending_connection).uid,s.(*pending_connection).rid,new_session);
					}
				case *KcpSession:
					if !s.(*KcpSession).congestion(){
						s.(*KcpSession).chPending<-p;
					}
				default:
					p.Return();
				}
			}else if len(evnt)<cap(evnt){
				evnt<-p;
			}else{
				p.Return();
			}
		}
	},ch_handshake);
}
func (this *KcpGateway)add_session(session *KcpSession){
	this.sessions.Store(session.RemoteAddr.String(),session);
}
func (this *KcpGateway)del_session(session *KcpSession){
	this.sessions.Delete(session.RemoteAddr.String());
}

func (this *KcpGateway)start(){
	this.once_start.Do(func() {
		this.start_main_proc();
	})
}

func StartGateway(addr string,on_session func(uid,rid uint32,session *KcpSession))(error){
	udpaddr, err := net.ResolveUDPAddr("udp", addr);
	if err!=nil{
		return err;
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err!=nil{
		return err;
	}
	ipv4.NewConn(conn).SetTOS(0<<2);
	conn.SetReadBuffer(32*utils.MaxPktSize*utils.MaxPktSize);
	conn.SetWriteBuffer(32*utils.MaxPktSize*utils.MaxPktSize);
	s:=new(KcpGateway);
	s.conn=conn;
	s.on_session=on_session;
	s.start();
	return nil;
}
