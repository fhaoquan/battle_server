package old

import (
	"github.com/xtaci/kcp-go"
	"time"
	"github.com/sirupsen/logrus"
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
	l,e:=kcp.ListenWithOptions(s.addr_string,nil,0,0);
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
					go func(c *kcp.UDPSession){
						d:=make([]byte,1024);
						for i:=0;i<1000;i++{
							if n,e:=c.Read(d);e==nil{
								logrus.Error("recved : ",d[:n])
							}
							time.Sleep(time.Second);
						}

						c.Close();
						logrus.Error("one session closed :",c.RemoteAddr())
					}(conn);
					//world.OnNewKCPConnection(conn);
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
