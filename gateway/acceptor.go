package gateway

import (
	"net"
	"../utils"
	"github.com/sirupsen/logrus"
	"encoding/binary"
	"sync"
)

var acceptor_start_once	=new(sync.Once);
var recv_map			=new(sync.Map);
var acceptor_buf		=make([]byte,1024);

func check_packet(pkt []byte)(bool){
	switch {
	case len(pkt)!=10:
		return false;
	case binary.LittleEndian.Uint16(pkt)!=54321:
		return false;
	default:
		return true;
	}
}

func create_session(uid uint32,rid uint32)(*session){
	receiver:=GetReceiver(rid,uid);
	if receiver==nil {
		return nil;
	}
	s,e:=new_session(rid,uid);
	if e!=nil{
		logrus.Error(e);
		return nil
	}
	receiver.NewSession(s);
	return s;
}
func report_session(conn *net.UDPConn,addr net.Addr,s *session){
	if s!=nil{
		logrus.Info("new connect rid=",s.rid,"uid=",s.uid,"local address=",s.conn.LocalAddr().String())
		binary.LittleEndian.PutUint16(acceptor_buf,10);
		binary.LittleEndian.PutUint32(acceptor_buf[2:],uint32(s.conn.LocalAddr().(*net.UDPAddr).Port));
		binary.LittleEndian.PutUint32(acceptor_buf[6:],uint32(s.conn.LocalAddr().(*net.UDPAddr).Port));
		conn.WriteTo(acceptor_buf[:10],addr);
	}
}

func accept_proc(conn *net.UDPConn){
	buf:=make([]byte,utils.MaxPktSize);
	for{
		l,a,e:=conn.ReadFrom(buf);
		switch{
		case e!=nil:
			logrus.Error("accept_proc ReadFrom",e)
		case !check_packet(buf[:l]):
			logrus.Error("accept_proc ReadFrom check pkt error");
		default:
			uid:=binary.LittleEndian.Uint32(buf[2:]);
			rid:=binary.LittleEndian.Uint32(buf[6:]);
			logrus.Info("recv connect request rid=",rid,"uid=",uid)
			report_session(
				conn,
				a,
				create_session(uid,rid),
				);
		}
	}
}
