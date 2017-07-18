package proc
import (
	"../udp"
	"../utils"
	"../battle"
	"net"
	"github.com/sirupsen/logrus"
	"time"
	"encoding/binary"
)

type udp_packet_carrier struct{
	len int;
	addr *net.UDPAddr;
	data []byte;
}
func (pkt *udp_packet_carrier) GetLen()int{
	return pkt.len;
}
func (pkt *udp_packet_carrier) GetAddr()*net.UDPAddr{
	return pkt.addr;
}
func (pkt *udp_packet_carrier) SetLen(v int){
	pkt.len=v;
}
func (pkt *udp_packet_carrier) SetAddr(v *net.UDPAddr){
	pkt.addr=v;
}
func (pkt *udp_packet_carrier) GetData()[]byte{
	return pkt.data;
}
func (pkt *udp_packet_carrier) GetPacketID()uint16{
	return binary.BigEndian.Uint16(pkt.GetData()[0:2]);
}

type i_msg_receiver interface {

}

type context struct {
	the_battle *battle.Battle;
	udp_connection *udp.UDPServer;
	udp_recv_memory_pool *utils.MemoryPool;
	udp_send_memory_pool *utils.MemoryPool;
	recv_chan chan *utils.CachedData;
	send_chan chan *utils.CachedData;
	tick_chan chan int;
}
func (c *context)udp_recv_proc(){
	for{
		empty_msg:=c.udp_recv_memory_pool.Get();
		err:=c.udp_connection.ReadMsg(empty_msg.GetUserData().(*udp_packet_carrier));
		if(err==nil){
			c.recv_chan<-empty_msg;
		}else{
			logrus.Info(err);
		}
	}
}
func (c *context)udp_send_proc(){
	for{
		select {
		case msg:=<- c.send_chan:
			err:=c.udp_connection.SendMsg(msg.GetUserData().(*udp_packet_carrier));
			msg.Return();
			if(err!=nil){
				logrus.Info(err);
			}
		}
	}
}
func (c *context)tick_proc(){
	i:=0;
	for{
		time.Sleep(time.Millisecond*30);
		c.tick_chan<-i;
		i++;
	}
}
func (c *context)main_proc(){
	for{
		select {
		case msg:=<-c.recv_chan:
			c.the_battle.OnMsg(msg.GetUserData().(*udp_packet_carrier));
			msg.Return();
		case <-c.tick_chan:
			c.the_battle.OnTick();
		}
	}
}
