package room

import (
	"../utils"
	"net"
)

type Player struct{
	flag int;
	id uint32;
	name string;
	udp_addr net.Addr;
	kcp_chan chan utils.IKcpResponse;
	udp_chan chan utils.IUdpResponse;
}
func (p *Player)SetKcpSender(c chan utils.IKcpResponse){
	p.kcp_chan=c;
}
func (p *Player)SetUDPSender(c chan utils.IUdpResponse){
	p.udp_chan=c;
}
func (p *Player)SetUDPAddr(addr net.Addr){
	p.udp_addr=addr;
}
func (p *Player)SendUDP(response utils.IUdpResponse){
	if p.udp_chan!=nil{
		response.SetAdr(p.udp_addr);
		p.udp_chan<-response;
	}
}
func (p *Player)SendKCP(response utils.IKcpResponse){
	if p.kcp_chan!=nil{
		p.kcp_chan<-response;
	}
}