package kcp_server

import (
	"github.com/xtaci/kcp-go"
	"time"
	"github.com/sirupsen/logrus"
	"../../world"
	"github.com/xtaci/smux"
)
type kcp_config interface {
	GetAddr()string;
}
type KcpServer struct {
	addr_string string;
	CloseSignal chan int;
}
func test1(c1 *kcp.UDPSession){
	d:=make([]byte,16);
	for{
		c1.Read(d);
		logrus.Error(string(d));
	}

}
func test2(c1 *kcp.UDPSession){
	d:=make([]byte,16);
	config:=smux.DefaultConfig();
	config.KeepAliveTimeout=time.Second*2;
	config.KeepAliveInterval=time.Second*1;
	session,_:=smux.Server(c1,config)
	stream,e:=session.AcceptStream();
	if(e!=nil){
		logrus.Error(e)
		//return ;
	}
	logrus.Error(session);
	j:=10;
	for i:=0;i<j;i++{
		logrus.Error("in session :",c1.RemoteAddr())
		n,e:=stream.Read(d);
		if e==nil{
			if(string(d[:n])=="1234"){
				j=10;
			}else{
				j=100;
			}
			stream.Write(d);
		}else{
			logrus.Error(e);
		}
	}
	logrus.Error("......",string(d));
	stream.Close();
	session.Close();
	logrus.Error("one session exited",c1.RemoteAddr())
}
func (s *KcpServer)StartAt(world *world.World)(error){
	l,e:=kcp.ListenWithOptions(s.addr_string,nil,10,3);
	if(e!=nil){
		return e;
	}
	if e=l.SetReadBuffer(1024*16);e!=nil{
		return e;
	}
	if e=l.SetWriteBuffer(1024*16);e!=nil{
		return e;
	}
	if e=l.SetDSCP(0);e!=nil{
		return e;
	}
	go func(){
		defer func(){
			close(s.CloseSignal);
			l.Close();
		}()
		for{
			select {
			case _,ok:=<-s.CloseSignal:
				if(ok){
					return;
				}
			default:
				l.SetDeadline(time.Now().Add(time.Second*1));
				conn,e:=l.AcceptKCP();
				if(e==nil){
					conn.SetReadBuffer(1024*10);
					conn.SetWriteBuffer(1024*10);
					conn.SetWindowSize(32, 32);
					conn.SetNoDelay(1, 5, 2, 1);
					conn.SetStreamMode(true);
					conn.SetMtu(1400);
					logrus.Error("one session connected :",conn.RemoteAddr())
					//go test1(conn);
					world.OnNewKCPConnection(conn);
				}else if err,ok:=e.(interface{Timeout()bool});!ok||!err.Timeout(){
					logrus.Error(e);
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
