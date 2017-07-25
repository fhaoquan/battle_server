package udp_session

import (
	"net"
	"fmt"
	"../../utils"
	"sync"
)

type session struct {
	m sync.RWMutex;
	con net.PacketConn;
	clients map[uint32]*net.UDPAddr;
}

func (s *session)for_each_client(f func(*net.UDPAddr)bool){
	defer s.m.RUnlock();
	s.m.RLock();
	for _,a:=range s.clients{
		if(!f(a)){
			return ;
		}
	}
}

func (s *session)get_client(uid uint32)*net.UDPAddr{
	defer s.m.RUnlock();
	s.m.RLock();
	return s.clients[uid];
}

func (s *session)set_client(uid uint32,addr *net.UDPAddr){
	s.m.Lock();
	s.clients[uid]=addr;
	s.m.Unlock();
}

func NewSession(port int)(*session,error){
	if adr,err:=net.ResolveUDPAddr("udp",fmt.Sprint(":",port));err!=nil{
		return nil,err;
	}else if con,err:=net.ListenUDP("udp", adr);err!=nil{
		return nil,err;
	}else{
		con.SetWriteBuffer(utils.MaxPktSize*16);
		con.SetReadBuffer(utils.MaxPktSize*16);
		return &session{
			con:con,
			clients:make(map[uint32]*net.UDPAddr),
		},nil;
	}
}
