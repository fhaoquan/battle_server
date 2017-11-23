package gateway

import (
	"net"
	"golang.org/x/net/ipv4"
	"../utils"
	"github.com/sirupsen/logrus"
)
type Session interface {
	RoomID()(uint32)
	UserID()(uint32)
	Start(k chan utils.I_REQ)
	Close()
	KcpSend(data []byte)(bool)
	UdpSend(data []byte)(bool)
}
type Receiver interface {
	NewSession(s Session);
	DelSession(s Session);
}

func AddReceiver(rid uint32,uid uint32,r Receiver){
	logrus.Error("AddReceiver",rid,uid);
	k:=uint64(uint64(rid)<<32)|uint64(uid)
	recv_map.Store(k,r);
}
func DelReceiver(rid uint32,uid uint32,){
	k:=uint64(uint64(rid)<<32)|uint64(uid)
	recv_map.Delete(k);
}
func GetReceiver(rid uint32,uid uint32)(Receiver){
	k:=uint64(uint64(rid)<<32)|uint64(uid)
	r,ok:=recv_map.Load(k);
	if ok{
		return r.(Receiver);
	}
	return nil;
}

func Start(addr string)(error){
	udpaddr, err := net.ResolveUDPAddr("udp", addr);
	if err!=nil{
		return err;
	}
	conn, err := net.ListenUDP("udp", udpaddr)
	if err!=nil{
		return err;
	}
	ipv4.NewConn(conn).SetTOS(0<<2);
	conn.SetReadBuffer(64*utils.MaxPktSize);
	conn.SetWriteBuffer(64*utils.MaxPktSize);
	acceptor_start_once.Do(func() {
		go accept_proc(conn);
	})
	return nil;
}
