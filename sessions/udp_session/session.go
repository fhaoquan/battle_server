package udp_session

import (
	"net"
	"fmt"
	"../../utils"
)

func TryListen(port int)(net.PacketConn,error){
	if adr,err:=net.ResolveUDPAddr("udp",fmt.Sprint(":",port));err!=nil{
		return nil,err;
	}else if con,err:=net.ListenUDP("udp", adr);err!=nil{
		return nil,err;
	}else{
		con.SetWriteBuffer(utils.MaxPktSize*16);
		con.SetReadBuffer(utils.MaxPktSize*16);
		return con,nil;
	}
}
