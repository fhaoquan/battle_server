package room

import (
	"../sessions/packet"
	"net"
)

type Player struct{
	flag int;
	id uint32;
	name string;
	udp_addr net.Addr;
	kcp_chan chan []byte;
	udp_chan chan func(func(addr net.Addr,data []byte));
}
func (p *Player)SetKcpSender(c chan []byte){
	p.kcp_chan=c;
}
func (p *Player)SetUDPSender(c chan func(func(addr net.Addr,data []byte))){
	p.udp_chan=c;
}
func (p *Player)SetUDPAddr(addr net.Addr){
	p.udp_addr=addr;
}
func (p *Player)SendUDP(response packet.IUdpResponse){
	p.udp_chan<-func(f func(addr net.Addr,data []byte)){
		f(p.udp_addr,response.GetBDY());
	}
}
func (p *Player)SendKCP(response packet.IKcpResponse){
	p.kcp_chan<-response.GetBDY();
}