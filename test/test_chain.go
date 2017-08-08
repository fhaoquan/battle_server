package test

import "net"

type kcp_packet struct {
	l uint16;
	u uint32;
	r uint32;
	b []byte;
}

func (me *kcp_packet)read(conn net.Conn)(error){
	return nil;
}

func (me *kcp_packet)cache(c chan *kcp_packet){
	c<-me;
}

func (me *kcp_packet)run(handlers []func([]byte)interface{})interface{}{
	return handlers[me.b[10]](me.b[11:]);
}

type kcp_chain struct {
	c chan *kcp_packet
}

func test(p *kcp_packet){
	p.read(nil);
	p.cache(nil);
}
