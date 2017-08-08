package kcp_session

import (
	"github.com/xtaci/kcp-go"
	"time"
	log "github.com/sirupsen/logrus"
	"../../world"
)
type kcp_config interface {
	GetAddr()string;
}
type KcpServer struct {
	addr_string string;
	CloseSignal chan int;
}
func (s *KcpServer)StartAt(world *world.World)(error){
	l,e:=kcp.Listen(s.addr_string);
	if(e!=nil){
		return e;
	}
	lis:=l.(*kcp.Listener);
	if e=lis.SetReadBuffer(1024*16);e!=nil{
		return e;
	}
	if e=lis.SetWriteBuffer(1024*16);e!=nil{
		return e;
	}
	if e=lis.SetDeadline(time.Now().Add(time.Second*5));e!=nil{
		return e;
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
					NewSession(conn).StartAt(world);
				}else if err,ok:=e.(interface{Timeout()bool});!ok||!err.Timeout(){
					log.Error(e);
				}
			}
		}
	}();
	return nil;
}
func NewKcpServer(addr string)(*KcpServer){
	s:=&KcpServer{
		addr_string:addr,
		CloseSignal:make(chan int,1),
	}
	return s;
}
