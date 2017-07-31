package room

type Player struct{
	flag int;
	id uint32;
	name string;
	kcp_chan chan []byte
}
func (p *Player)SetKcpSender(c chan []byte){
	p.kcp_chan=c;
}