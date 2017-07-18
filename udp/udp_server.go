package udp

import "net"

type i_udp_msg_carrier interface {
	GetData()[]byte;
	SetLen(v int);
	GetLen()int;
	SetAddr(v *net.UDPAddr);
	GetAddr()*net.UDPAddr;
}

type UDPServer struct {
	conn *net.UDPConn;
	ip string;
	port int;
}

func (server *UDPServer)ReadMsg(msg i_udp_msg_carrier)(error){
	len, remoteAddr, err :=server.conn.ReadFromUDP(msg.GetData());
	msg.SetLen(len);
	msg.SetAddr(remoteAddr);
	return err;
}
func (server *UDPServer)SendMsg(msg i_udp_msg_carrier)(error){
	_,err:=server.conn.WriteTo(msg.GetData()[0:msg.GetLen()],msg.GetAddr());
	return err;
}

func NewUDPServer(addr string)(*UDPServer,error){
	u_addr,err:=net.ResolveUDPAddr("udp",addr);
	if(err!=nil){
		return nil,err;
	}
	conn, err :=net.ListenUDP("udp", u_addr);
	if(err==nil){
		return nil,err;
	}
	res:=&UDPServer{
		conn:conn,
	}
	return res,nil;
}

