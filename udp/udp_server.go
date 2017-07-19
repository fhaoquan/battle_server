package udp

import (
	"net"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"../utils"
	"io"
	"encoding/binary"
)

type i_udp_msg_carrier interface {
	GetData()[]byte;
	SetLen(v int);
	GetLen()int;
	SetAddr(v *net.UDPAddr);
	GetAddr()*net.UDPAddr;
}
type i_udp_server_owner interface {
	GetID();
	OnMsg([]byte,int);
}
type UDPServer struct {
	conn *net.UDPConn;
	ip string;
	port int;
	owner i_udp_server_owner;
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
func (server *UDPServer)recv_proc(){
	err:=func()(e error){
		defer func(){
			if err:=recover();err!=nil{
				e=errors.New(fmt.Sprint(err));
			}
		}();
		d:=make([]byte,utils.MaxPktSize)
		for{
			io.ReadFull(server.conn,d[0:2]);
			if _,err:=io.ReadFull(server.conn,d[0:2]);err!=nil{
				return err;
			}
			l:=binary.BigEndian.Uint16(d[0:2]);
			if _,err:=io.ReadFull(server.conn,d[2:l+2]);err!=nil{
				return  err;
			}
			//server.owner.OnMsg()

		}
	}();
	if(err!=nil){
		logrus.Error(err);
	}
}
func (server *UDPServer)Start(){

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

