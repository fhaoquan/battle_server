package udp_server

import (
	"net"
	"sync"
	"github.com/sirupsen/logrus"
	"../../utils"
)


var UdpSlot = make([]chan *utils.UdpReq,1024);

type UdpGateway struct {
	once_start			sync.Once;
	conn				*net.UDPConn;
}
func (this *UdpGateway)go_gateway_kernel_proc(f func(a interface{}),a interface{}){
	go func(){
		defer func() {
			if e:=recover();e!=nil{
				logrus.Fatal(e);
			}
		}()
		f(a);
	}();
}
func (this *UdpGateway)start_main_proc(){
	this.go_gateway_kernel_proc(func(a interface{}) {
		pool:=utils.NewMemoryPool(64*1000, func(impl utils.ICachedData) utils.ICachedData {
			return &utils.UdpReq{
				nil,&utils.KcpReq{
					impl,make([]byte,utils.MaxPktSize),
				},
			}
		})
		e:=error(nil);
		for{
			p:=pool.Pop().(*utils.UdpReq);
			_,p.ADR,e=this.conn.ReadFrom(p.Data);
			if e!=nil{
				p.Return();
				logrus.Error(e);
				continue;
			}
			rid:=p.GetRID();
			if rid<10000{
				p.Return();
				continue;
			}
			if int(rid-10000)>=len(UdpSlot){
				p.Return();
				continue;
			}
			c:=UdpSlot[rid-10000];
			if c==nil{
				p.Return();
				continue;
			}
			if cap(c)<=len(c){
				p.Return();
				continue;
			}
			c<-p;
		}
	},nil);
}
func (this *UdpGateway)start(){
	this.once_start.Do(func() {
		this.start_main_proc();
	})
}
func StartGateway(addr string){
	udpaddr, err := net.ResolveUDPAddr("udp", addr);
	if err!=nil{
		logrus.Fatal(err);
		return;
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err!=nil{
		logrus.Fatal(err);
		return;
	}
	conn.SetReadBuffer(32*1024*1024);
	conn.SetWriteBuffer(32*1024*1024);
	s:=new(UdpGateway);
	s.conn=conn;
	s.start();
	return;
}