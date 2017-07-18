package udp

import (
	"github.com/xtaci/kcp-go"
	"time"
	log "github.com/sirupsen/logrus"
	"net"
)
type kcp_config interface {
	GetAddr()string;
}
type KcpServer struct {
	CloseSignal chan int;
}

func StartNewKcpServer(addr string,on_connection func(conn net.Conn))(*KcpServer,error){
	l,e:=kcp.Listen(addr);
	if(e!=nil){
		return nil,e;
	}
	lis:=l.(*kcp.Listener);
	if e=lis.SetReadBuffer(1024*16);e!=nil{
		return nil,e;
	}
	if e=lis.SetWriteBuffer(1024*16);e!=nil{
		return nil,e;
	}
	if e=lis.SetDeadline(time.Now().Add(time.Second*5));e!=nil{
		return nil,e;
	}
	s:=&KcpServer{
		CloseSignal:make(chan int,1),
	}
	go func(){
		defer func(){
			close(s.CloseSignal);
			lis.Close();
		}()
		for{
			select {
			case _,ok:=<-s.CloseSignal:
				if(ok){
					return;
				}
			default:
				lis.SetDeadline(time.Now().Add(time.Second*1));
				conn,e:=lis.AcceptKCP();
				if(e==nil){
					conn.SetReadBuffer(1024*10);
					conn.SetWriteBuffer(1024*10);
					conn.SetWindowSize(32, 32);
					conn.SetNoDelay(1, 5, 2, 1);
					conn.SetStreamMode(true);
					conn.SetMtu(1400);
					on_connection(conn);
				}else if err,ok:=e.(interface{Timeout()bool});!ok||!err.Timeout(){
					log.Error(e);
				}
			}
		}
	}();

	return s,nil;
}
